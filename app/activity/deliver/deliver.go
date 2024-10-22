package deliver

// TBD replace all instances of ap.Activity with activitypub.Activity ?
// This would allow to get rid of all the PostId stuff and move this in
// app/activity/deliver

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
	"github.com/sfomuseum/go-activitypub/database"
	"github.com/sfomuseum/go-activitypub/followers"
	// "github.com/sfomuseum/go-activitypub/posts"
	"github.com/sfomuseum/go-activitypub/queue"
	"github.com/sfomuseum/go-pubsub/subscriber"
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

	// logger := slog.Default()

	accounts_db, err := database.NewAccountsDatabase(ctx, opts.AccountsDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create accounts database, %w", err)
	}

	defer accounts_db.Close(ctx)

	activities_db, err := database.NewActivitiesDatabase(ctx, opts.ActivitiesDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create activities database, %w", err)
	}

	defer activities_db.Close(ctx)

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

	deliverActivityTo := func(ctx context.Context, activity_id int64, recipient string) error {

		logger := slog.Default()

		logger = logger.With("activity id", activity_id)
		logger = logger.With("to", recipient)

		activity, err := activities_db.GetActivityWithId(ctx, activity_id)

		if err != nil {
			logger.Error("Failed to retrieve activity", "error", err)
			return fmt.Errorf("Failed to retrieve activity, %w", err)
		}

		acct, err := accounts_db.GetAccountWithId(ctx, activity.AccountId)

		if err != nil {
			logger.Error("Failed to retrieve account", "account id", activity.AccountId)
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

		// START OF wrangle mentions if activity is post

		switch activity.ActivityType {
		case activitypub.PostActivityType:

			post_id := activity.ActivityTypeId
			logger = logger.With("post id", post_id)

			mentions := make([]*activitypub.PostTag, 0)

			mentions_cb := func(ctx context.Context, t *activitypub.PostTag) error {
				mentions = append(mentions, t)
				return nil
			}

			err = post_tags_db.GetPostTagsForPost(ctx, post_id, mentions_cb)

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

		case activitypub.BoostActivityType:

			// TBD... do this for all activities not just boosts?

			if !is_allowed {

				ap_activity, err := activity.UnmarshalActivity()

				if err != nil {
					logger.Error("Failed to unmarshal activity", "error", err)
					return fmt.Errorf("Failed to unmarshal activity, %w", err)
				}

				for _, addr := range ap_activity.Cc {

					if addr == recipient {
						logger.Info("Recipient is not allowed/followed but is listed in Cc", "addr", addr)
						is_allowed = true
						break
					}
				}
			}

		default:
			// pass
		}

		// END OF wrangle mentions if activity is post

		if !is_allowed {
			logger.Error("Recipient is flagged as 'not allowed' to have message delivered")
			return nil
		}

		deliver_opts := &queue.DeliverActivityOptions{
			To:                 recipient,
			Activity:           activity,
			AccountsDatabase:   accounts_db,
			DeliveriesDatabase: deliveries_db,
			URIs:               opts.URIs,
			MaxAttempts:        opts.MaxAttempts,
		}

		logger.Debug("Deliver activity")

		err = queue.DeliverActivity(ctx, deliver_opts)

		if err != nil {
			return fmt.Errorf("Failed to deliver activity, %w", err)
		}

		logger.Info("Activity delivered")
		return nil
	}

	// END OF...

	switch opts.Mode {
	case "cli":

		// This needs to be refactored to be more like the "lambda" and "pubsub" modes below.
		// Or maybe there needs to be a second "cli-single" (or whatever) mode that operates
		// on delivering a post to a single recipient. As it is the "cli" mode schedules delivery
		// of a post to all the people following the post author.

		/*
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

		*/

	case "lambda":

		// For processing posts scheduled to be delivered to individual recipients  an AWS SQS queue

		handler := func(ctx context.Context, sqsEvent events.SQSEvent) error {

			for _, message := range sqsEvent.Records {

				logger := slog.Default()
				logger = logger.With("message (sqs) id", message.MessageId)

				var ps_opts *queue.PubSubDeliveryQueueOptions

				err := json.Unmarshal([]byte(message.Body), &ps_opts)

				if err != nil {
					logger.Error("Failed to unmarshal deliver options", "error", err)
					return fmt.Errorf("Failed to unmarshal deliver options, %w", err)
				}

				return deliverActivityTo(ctx, ps_opts.ActivityId, ps_opts.To)
			}

			return nil
		}

		lambda.Start(handler)
		return nil

	case "pubsub":

		// For processing posts in a shared "pubsub" environment, for example Redis. Uses sfomuseum/go-pubsub
		// to send and receive messages. This is mostly for being able to debug the "lambda" handler in local
		// dev setup.

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

						err := deliverActivityTo(ctx, ps_opts.ActivityId, ps_opts.To)

						if err != nil {
							logger.Error("Failed to deliver activity", "activity id", ps_opts.ActivityId, "to", ps_opts.To, "error", err)
						}
					}

				default:
					//
				}
			}
		}()

		logger.Info("Listening for messages on pubsub channel")
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
