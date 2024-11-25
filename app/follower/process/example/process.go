package example

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sfomuseum/go-activitypub/database"
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

	messages_db, err := database.NewMessagesDatabase(ctx, opts.MessagesDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create new messages database, %w", err)
	}

	defer messages_db.Close(ctx)

	notes_db, err := database.NewNotesDatabase(ctx, opts.NotesDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create new notes database, %w", err)
	}

	defer notes_db.Close(ctx)

	accounts_db, err := database.NewAccountsDatabase(ctx, opts.AccountsDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create new accounts database, %w", err)
	}

	defer accounts_db.Close(ctx)

	properties_db, err := database.NewPropertiesDatabase(ctx, opts.PropertiesDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create new properties database, %w", err)
	}

	defer properties_db.Close(ctx)

	activities_db, err := database.NewActivitiesDatabase(ctx, opts.ActivitiesDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create new activities database, %w", err)
	}

	defer activities_db.Close(ctx)

	posts_db, err := database.NewPostsDatabase(ctx, opts.PostsDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create new posts database, %w", err)
	}

	defer posts_db.Close(ctx)

	post_tags_db, err := database.NewPostTagsDatabase(ctx, opts.PostTagsDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create new post tags database, %w", err)
	}

	defer post_tags_db.Close(ctx)

	deliveries_db, err := database.NewDeliveriesDatabase(ctx, opts.DeliveriesDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create new deliveries database, %w", err)
	}

	defer deliveries_db.Close(ctx)

	followers_db, err := database.NewFollowersDatabase(ctx, opts.FollowersDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create new followers database, %w", err)
	}

	defer followers_db.Close(ctx)

	_, err = queue.NewDeliveryQueue(ctx, opts.DeliveryQueueURI)

	if err != nil {
		return fmt.Errorf("Failed to create new delivery queue, %w", err)
	}

	// Note: Don't close delivery_q in a Lambda context. This will trigger errors like this:
	// "Failed to send message, pubsub: Topic has been Shutdown (code=FailedPrecondition)"

	process_follower := func(ctx context.Context, follower_id int64) error {

		logger := slog.Default()
		logger = logger.With("method", "process_follow")
		logger = logger.With("follower id", follower_id)

		logger.Info("Process follow")

		// Message is the thing which was dispatched by the process message queue (in www/inbox_post.go)

		f, err := followers_db.GetFollowerWithId(ctx, follower_id)

		if err != nil {
			logger.Error("Failed to retrieve message", "error", err)
			return fmt.Errorf("Failed to retrieve message, %w", err)
		}

		// Account is the account that the note (message) was delivered to

		msg_acct, err := accounts_db.GetAccountWithId(ctx, f.AccountId)

		if err != nil {
			logger.Error("Failed to retrieve account for follow", "account id", f.AccountId, "error", err)
			return fmt.Errorf("Failed to retrieve account for message, %w", err)
		}

		logger = logger.With("follow account id", msg_acct.Id)
		logger = logger.With("follow follower address", f.FollowerAddress)

		logger.Info("Your code goes here (parse note, etc.)")
		return nil
	}

	// Process multiple follower IDs
	// Maybe try to do this concurrently?

	process_followers := func(ctx context.Context, follower_ids ...int64) error {

		for _, id := range follower_ids {

			err := process_follower(ctx, id)

			if err != nil {
				slog.Error("Failed to process follower", "id", id, "error", err)
			}
		}

		return nil
	}

	// Actually start the application

	switch opts.Mode {
	case "cli":

		return process_followers(ctx, opts.FollowerIds...)
	case "lambda":

		handler := func(ctx context.Context, sqsEvent events.SQSEvent) error {

			follower_ids := make([]int64, len(sqsEvent.Records))

			for idx, message := range sqsEvent.Records {

				logger := slog.Default()
				logger = logger.With("message id", message.MessageId)

				// logger.Debug("SQS", "message", message.Body)

				var follower_id int64

				err := json.Unmarshal([]byte(message.Body), &follower_id)

				if err != nil {
					logger.Error("Failed to unmarshal message", "error", err)
					return fmt.Errorf("Failed to unmarshal message ID, %w", err)
				}

				follower_ids[idx] = follower_id
			}

			return process_followers(ctx, follower_ids...)
		}

		lambda.Start(handler)

	default:
		return fmt.Errorf("Invalid or unsupported mode")
	}

	return nil
}
