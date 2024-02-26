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

	deliveries_db, err := activitypub.NewDeliveriesDatabase(ctx, opts.DeliveriesDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create instantiate deliveries database, %w", err)
	}

	defer deliveries_db.Close(ctx)

	delivery_q, err := activitypub.NewDeliveryQueue(ctx, opts.DeliveryQueueURI)

	if err != nil {
		return fmt.Errorf("Failed to create new delivery queue, %w", err)
	}

	switch opts.Mode {
	case "cli":

		accounts_db, err := activitypub.NewAccountsDatabase(ctx, opts.AccountsDatabaseURI)

		if err != nil {
			return fmt.Errorf("Failed to create new database, %w", err)
		}

		defer accounts_db.Close(ctx)

		followers_db, err := activitypub.NewFollowersDatabase(ctx, opts.FollowersDatabaseURI)

		if err != nil {
			return fmt.Errorf("Failed to instantiate followers database, %w", err)
		}

		defer followers_db.Close(ctx)

		posts_db, err := activitypub.NewPostsDatabase(ctx, opts.PostsDatabaseURI)

		if err != nil {
			return fmt.Errorf("Failed to create instantiate posts database, %w", err)
		}

		defer posts_db.Close(ctx)

		post, err := posts_db.GetPostWithId(ctx, opts.PostId)

		if err != nil {
			return fmt.Errorf("Failed to retrieve post, %w", err)
		}

		deliver_opts := &activitypub.DeliverPostToFollowersOptions{
			AccountsDatabase:   accounts_db,
			FollowersDatabase:  followers_db,
			DeliveriesDatabase: deliveries_db,
			DeliveryQueue:      delivery_q,
			Post:               post,
			URIs:               opts.URIs,
		}

		err = activitypub.DeliverPostToFollowers(ctx, deliver_opts)

		if err != nil {
			return fmt.Errorf("Failed to deliver post, %w", err)
		}

		return nil

	case "lambda":

		handler := func(ctx context.Context, snsEvent events.SNSEvent) error {

			for _, record := range snsEvent.Records {

				var opts *activitypub.DeliverPostOptions
				opts.DeliveriesDatabase = deliveries_db // It's not great having to do this so what is better...?

				err := json.Unmarshal([]byte(record.SNS.Message), &opts)

				if err != nil {
					slog.Error("Failed to unmarshal post options", "error", err)
					return fmt.Errorf("Failed to unmarshal post options, %w", err)
				}

				err = activitypub.DeliverPost(ctx, opts)

				if err != nil {
					slog.Error("Failed to deliver post", "error", err)
					return fmt.Errorf("Failed to deliver post, %w", err)
				}

			}

			return nil
		}

		lambda.Start(handler)
		return nil

	default:
		return fmt.Errorf("Invalid or unsupported mode")
	}
}
