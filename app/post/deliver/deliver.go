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
	ap_slog "github.com/sfomuseum/go-activitypub/slog"
)

func Run(ctx context.Context, logger *slog.Logger) error {
	fs := DefaultFlagSet()
	return RunWithFlagSet(ctx, fs, logger)
}

func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet, logger *slog.Logger) error {

	opts, err := OptionsFromFlagSet(ctx, fs)

	if err != nil {
		return fmt.Errorf("Failed to derive options from flagset, %w", err)
	}

	return RunWithOptions(ctx, opts, logger)
}

func RunWithOptions(ctx context.Context, opts *RunOptions, logger *slog.Logger) error {

	ap_slog.ConfigureLogger(logger, opts.Verbose)

	accounts_db, err := activitypub.NewAccountsDatabase(ctx, opts.AccountsDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create new database, %w", err)
	}

	defer accounts_db.Close(ctx)

	deliveries_db, err := activitypub.NewDeliveriesDatabase(ctx, opts.DeliveriesDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create instantiate deliveries database, %w", err)
	}

	defer deliveries_db.Close(ctx)

	posts_db, err := activitypub.NewPostsDatabase(ctx, opts.PostsDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create instantiate posts database, %w", err)
	}

	defer posts_db.Close(ctx)

	post_tags_db, err := activitypub.NewPostTagsDatabase(ctx, opts.PostTagsDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create instantiate post tags database, %w", err)
	}

	defer post_tags_db.Close(ctx)

	followers_db, err := activitypub.NewFollowersDatabase(ctx, opts.FollowersDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to instantiate followers database, %w", err)
	}

	defer followers_db.Close(ctx)

	switch opts.Mode {
	case "cli":

		delivery_q, err := activitypub.NewDeliveryQueue(ctx, opts.DeliveryQueueURI)

		if err != nil {
			return fmt.Errorf("Failed to create new delivery queue, %w", err)
		}

		post, err := posts_db.GetPostWithId(ctx, opts.PostId)

		if err != nil {
			return fmt.Errorf("Failed to retrieve post, %w", err)
		}

		deliver_opts := &activitypub.DeliverPostToFollowersOptions{
			AccountsDatabase:   accounts_db,
			FollowersDatabase:  followers_db,
			PostTagsDatabase:   post_tags_db,
			DeliveriesDatabase: deliveries_db,
			DeliveryQueue:      delivery_q,
			Post:               post,
			URIs:               opts.URIs,
			MaxAttempts:        opts.MaxAttempts,
		}

		err = activitypub.DeliverPostToFollowers(ctx, deliver_opts)

		if err != nil {
			return fmt.Errorf("Failed to deliver post, %w", err)
		}

		return nil

	case "lambda":

		handler := func(ctx context.Context, sqsEvent events.SQSEvent) error {

			for _, message := range sqsEvent.Records {

				logger := slog.Default()
				logger = logger.With("message id", message.MessageId)

				// logger.Debug("SQS", "message", message.Body)

				var ps_opts *activitypub.PubSubDeliveryQueuePostOptions

				err := json.Unmarshal([]byte(message.Body), &ps_opts)

				if err != nil {
					logger.Error("Failed to unmarshal post options", "error", err)
					return fmt.Errorf("Failed to unmarshal post options, %w", err)
				}

				acct, err := accounts_db.GetAccountWithId(ctx, ps_opts.AccountId)

				if err != nil {
					slog.Error("Failed to retrieve account", "account id", ps_opts.Recipient, "error", err)
					return fmt.Errorf("Failed to retrieve account, %w", err)
				}

				logger = logger.With("account id", acct.Id)
				logger = logger.With("recipient", ps_opts.Recipient)

				is_follower, _, err := activitypub.IsFollower(ctx, followers_db, acct.Id, ps_opts.Recipient)

				if err != nil {
					slog.Error("Unable to determine if recipient is not following account", "recipient", ps_opts.Recipient, "error", err)
					return fmt.Errorf("Unable to determine if recipient is following account")
				}

				logger.Info("Follower check", "is_follower", is_follower)
				is_allowed := is_follower

				// We check this below after we've ensured recipient isn't mentioned in post (tags)

				post, err := posts_db.GetPostWithId(ctx, ps_opts.PostId)

				if err != nil {
					slog.Error("Failed to retrieve post", "post id", ps_opts.PostId, "error", err)
					return fmt.Errorf("Failed to retrieve post, %w", err)
				}

				logger = logger.With("post id", post.Id)

				if post.AccountId != acct.Id {
					slog.Error("Post owned by different account", "post account id", post.AccountId)
					return fmt.Errorf("Post owned by different account")
				}

				post_tags := make([]*activitypub.PostTag, 0)

				post_tags_cb := func(ctx context.Context, t *activitypub.PostTag) error {
					post_tags = append(post_tags, t)
					return nil
				}

				err = post_tags_db.GetPostTagsForPost(ctx, post.Id, post_tags_cb)

				if err != nil {
					slog.Error("Failed to retrieve post tags for post", "error", err)
					return fmt.Errorf("Failed to retrieve post tags for post, %w", err)
				}

				if !is_allowed && opts.AllowMentions && len(post_tags) > 0 {

					logger.Debug("Check to see whether recipient is listed in post tags", "count tags", len(post_tags))

					r_actor, err := activitypub.RetrieveActor(ctx, ps_opts.Recipient, opts.URIs.Insecure)

					if err != nil {
						logger.Warn("Failed to retrieve actor record for recipient", "recipient", ps_opts.Recipient, "error", err)
					} else {

						for _, t := range post_tags {

							// The old way
							// if t.Href == r_actor.Id {
								
							// The maybe? way
							// https://github.com/sfomuseum/go-activitypub/issues/3
							
							if t.Href == r_actor.URL {
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

				opts := &activitypub.DeliverPostOptions{
					From:               acct,
					To:                 ps_opts.Recipient,
					Post:               post,
					PostTags:           post_tags,
					URIs:               opts.URIs,
					DeliveriesDatabase: deliveries_db,
					MaxAttempts:        opts.MaxAttempts,
				}

				err = activitypub.DeliverPost(ctx, opts)

				if err != nil {
					slog.Error("Failed to deliver post", "message id", message.MessageId, "error", err)
					return fmt.Errorf("Failed to deliver post, %w", err)
				}

				logger.Info("Post delivered")
			}

			return nil
		}

		lambda.Start(handler)
		return nil

	default:
		return fmt.Errorf("Invalid or unsupported mode")
	}
}
