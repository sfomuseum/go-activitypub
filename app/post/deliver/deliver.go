package deliver

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sfomuseum/go-activitypub"
	"github.com/sfomuseum/go-activitypub/ap"
	"github.com/sfomuseum/go-pubsub/subscriber"	
	"github.com/sfomuseum/go-activitypub/database"
	"github.com/sfomuseum/go-activitypub/followers"
	"github.com/sfomuseum/go-activitypub/posts"
	"github.com/sfomuseum/go-activitypub/queue"
)

func Run(ctx context.Context) error {
	fs := DefaultFlagSet()
	return RunWithFlagSet(ctx, fs)
}

func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet) error {

	opts, err := OptionsFromFlagSet(ctx, fs)

	if err != nil {
		return fmt.Errorf("Failed to derive options from flagset, %w", err)
	}

	return RunWithOptions(ctx, opts)
}

func RunWithOptions(ctx context.Context, opts *RunOptions) error {

	if opts.Verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Verbose logging enabled")
	}

	logger := slog.Default()

	accounts_db, err := database.NewAccountsDatabase(ctx, opts.AccountsDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create new database, %w", err)
	}

	defer accounts_db.Close(ctx)

	deliveries_db, err := database.NewDeliveriesDatabase(ctx, opts.DeliveriesDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create instantiate deliveries database, %w", err)
	}

	defer deliveries_db.Close(ctx)

	posts_db, err := database.NewPostsDatabase(ctx, opts.PostsDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create instantiate posts database, %w", err)
	}

	defer posts_db.Close(ctx)

	post_tags_db, err := database.NewPostTagsDatabase(ctx, opts.PostTagsDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create instantiate post tags database, %w", err)
	}

	defer post_tags_db.Close(ctx)

	followers_db, err := database.NewFollowersDatabase(ctx, opts.FollowersDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to instantiate followers database, %w", err)
	}

	defer followers_db.Close(ctx)

	// START OF...
	
	deliverPostTo := func(ctx context.Context, post_id int64, recipient string) error {

		logger := slog.Default()
		
		logger = logger.With("to", recipient)
		logger = logger.With("post id", post_id)
		
		post, err := posts_db.GetPostWithId(ctx, post_id)
		
		if err != nil {
			logger.Error("Failed to retrieve post", "error", err)
			return fmt.Errorf("Failed to retrieve post, %w", err)
		}
		
		acct, err := accounts_db.GetAccountWithId(ctx, post.AccountId)
		
		if err != nil {
			logger.Error("Failed to retrieve account", "account id", post.AccountId)
			return fmt.Errorf("Failed to retrieve account, %w", err)
		}
		
		logger = logger.With("account id", acct.Id)
		
		is_follower, _, err := followers.IsFollower(ctx, followers_db, acct.Id, recipient)
		
		if err != nil {
			logger.Error("Unable to determine if recipient is not following account", "error", err)
			return fmt.Errorf("Unable to determine if recipient is following account")
		}
		
		logger.Info("Follower check", "is_follower", is_follower)
		is_allowed := is_follower

		mentions := make([]*activitypub.PostTag, 0)
		
		mentions_cb := func(ctx context.Context, t *activitypub.PostTag) error {
			mentions = append(mentions, t)
			return nil
		}
		
		err = post_tags_db.GetPostTagsForPost(ctx, post.Id, mentions_cb)
		
		if err != nil {
			slog.Error("Failed to retrieve post tags for post", "error", err)
			return fmt.Errorf("Failed to retrieve post tags for post, %w", err)
		}
		
		if !is_allowed && opts.AllowMentions && len(mentions) > 0 {
			
			logger.Debug("Check to see whether recipient is listed in post tags", "count tags", len(mentions))
			
			r_actor, err := ap.RetrieveActor(ctx, recipient, opts.URIs.Insecure)
			
			if err != nil {
				logger.Warn("Failed to retrieve actor record for recipient", "error", err)
			} else {
				
				for _, t := range mentions {
					
					// https://github.com/sfomuseum/go-activitypub/issues/3
					// if t.Href == r_actor.URL {
					
					// And yet it appears to actually be {ACTOR}.id however this
					// does not work (where "work" means open profile tab) in Ivory
					// yet because... I have no idea
					if t.Href == r_actor.Id {
						logger.Info("Recipient is included in post tags, allow delivery")
						is_allowed = true
						break
					}
				}
			}
			
		}
		
		if !is_allowed {
			logger.Error("Recipient is flagged as 'not allowed' to have message delivered")
			return nil
		}
		
		logger = logger.With("post id", post.Id)
		
		activity, err := posts.ActivityFromPost(ctx, opts.URIs, acct, post, mentions)
		
		if err != nil {
			return fmt.Errorf("Failed to create new (create) activity, %w", err)
		}
		
		logger = logger.With("activity id", activity.Id)
		
		deliver_opts := &queue.DeliverActivityOptions{
			To:                 recipient,
			Activity:           activity,
			PostId:             post.Id,
			AccountsDatabase:   accounts_db,
			DeliveriesDatabase: deliveries_db,
			URIs:               opts.URIs,
			MaxAttempts:        opts.MaxAttempts,
		}
		
		logger.Debug("Deliver activity")
		
		err = queue.DeliverActivity(ctx, deliver_opts)
		
		if err != nil {
			return fmt.Errorf("Failed to deliver post, %w", err)
		}
		
		logger.Info("Delivered post", "post url", acct.PostURL(ctx, opts.URIs, post).String())
		return nil
	}

	// END OF...
	
	switch opts.Mode {
	case "cli":

		delivery_q, err := queue.NewDeliveryQueue(ctx, opts.DeliveryQueueURI)

		if err != nil {
			return fmt.Errorf("Failed to create new delivery queue, %w", err)
		}

		post, err := posts_db.GetPostWithId(ctx, opts.PostId)

		if err != nil {
			return fmt.Errorf("Failed to retrieve post, %w", err)
		}

		acct, err := accounts_db.GetAccountWithId(ctx, post.AccountId)

		if err != nil {
			return fmt.Errorf("Failed to retrieve account %d, %w", post.AccountId, err)
		}

		mentions := make([]*activitypub.PostTag, 0)

		mentions_cb := func(ctx context.Context, t *activitypub.PostTag) error {
			mentions = append(mentions, t)
			return nil
		}

		err = post_tags_db.GetPostTagsForPost(ctx, post.Id, mentions_cb)

		if err != nil {
			return fmt.Errorf("Failed to retrieve tags for post, %w", err)
		}

		logger = logger.With("post id", post.Id)

		activity, err := posts.ActivityFromPost(ctx, opts.URIs, acct, post, mentions)

		if err != nil {
			return fmt.Errorf("Failed to create new (create) activity, %w", err)
		}

		logger = logger.With("activity id", activity.Id)

		deliver_opts := &queue.DeliverActivityToFollowersOptions{
			AccountsDatabase:   accounts_db,
			FollowersDatabase:  followers_db,
			DeliveriesDatabase: deliveries_db,
			DeliveryQueue:      delivery_q,
			Activity:           activity,
			PostId:             post.Id,
			Mentions:           mentions,
			URIs:               opts.URIs,
			MaxAttempts:        opts.MaxAttempts,
		}

		logger.Debug("Deliver activity")

		err = queue.DeliverActivityToFollowers(ctx, deliver_opts)

		if err != nil {
			return fmt.Errorf("Failed to deliver post, %w", err)
		}

		logger.Info("Delivered post", "post url", acct.PostURL(ctx, opts.URIs, post).String())
		return nil
		
	case "lambda":

		handler := func(ctx context.Context, sqsEvent events.SQSEvent) error {

			for _, message := range sqsEvent.Records {

				logger := slog.Default()
				logger = logger.With("message (sqs) id", message.MessageId)

				var ps_opts *queue.PubSubDeliveryQueueOptions

				err := json.Unmarshal([]byte(message.Body), &ps_opts)

				if err != nil {
					logger.Error("Failed to unmarshal post options", "error", err)
					return fmt.Errorf("Failed to unmarshal post options, %w", err)
				}

				return deliverPostTo(ctx, ps_opts.PostId, ps_opts.To)
			}

			return nil
		}

		lambda.Start(handler)
		return nil

	case "pubsub":

		logger := slog.Default()
		
		sub, err := subscriber.NewSubscriber(ctx, opts.SubscriberURI)

		if err != nil {
			return fmt.Errorf("Failed to create new subscriber, %w", err)
		}

		defer sub.Close()

		msg_ch := make(chan string)
		done_ch := make(chan bool)
		
		go func() {
			
			for {
				select {
				case <-ctx.Done():
					return
				case <-done_ch:
					return
				case msg := <-msg_ch:
					
					var ps_opts *queue.PubSubDeliveryQueueOptions
					
					err := json.Unmarshal([]byte(msg), &ps_opts)
					
					if err != nil {
						logger.Error("Failed to unmarshal post options", "error", err)
					} else {

						err := deliverPostTo(ctx, ps_opts.PostId, ps_opts.To)

						if err != nil {
							logger.Error("Failed to deliver post", "post id", ps_opts.PostId, "to", ps_opts.To, "error", err)
						}
					}
					
				default:
					//
				}
			}
		}()
		
		logger.Info("Listening for messages")
		err = sub.Listen(ctx, msg_ch)
		
		done_ch <- true
		
		if err != nil {
			return fmt.Errorf("Failed to listen, %v", err)
		}
		
	default:
		return fmt.Errorf("Invalid or unsupported mode")
	}

	return nil
}
