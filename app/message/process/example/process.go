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

	// Process a single message ID

	process_message := func(ctx context.Context, message_id int64) error {

		logger := slog.Default()
		logger = logger.With("method", "process_message")
		logger = logger.With("message id", message_id)

		logger.Info("Process message")

		// Message is the thing which was dispatched by the process message queue (in www/inbox_post.go)

		m, err := messages_db.GetMessageWithId(ctx, message_id)

		if err != nil {
			logger.Error("Failed to retrieve message", "error", err)
			return fmt.Errorf("Failed to retrieve message, %w", err)
		}

		// Note is the body (the actual AP note) that the message is associated with

		n, err := notes_db.GetNoteWithId(ctx, m.NoteId)

		if err != nil {
			logger.Error("Failed to retrieve note for message", "error", err)
			return fmt.Errorf("Failed to retrieve note for message, %w", err)
		}

		logger = logger.With("note id", n.Id)

		// Account is the account that the note (message) was delivered to

		msg_acct, err := accounts_db.GetAccountWithId(ctx, m.AccountId)

		if err != nil {
			logger.Error("Failed to retrieve account for message", "account id", m.AccountId, "error", err)
			return fmt.Errorf("Failed to retrieve account for message, %w", err)
		}

		logger = logger.With("message account id", msg_acct.Id)
		logger = logger.With("message account address", msg_acct.Address(opts.URIs.Hostname))

		logger.Info("Your code goes here (parse note, etc.)")
		return nil
	}

	// Process multiple message IDs
	// Maybe try to do this concurrently?

	process_messages := func(ctx context.Context, message_ids ...int64) error {

		for _, id := range message_ids {

			err := process_message(ctx, id)

			if err != nil {
				slog.Error("Failed to process message", "id", id, "error", err)
			}
		}

		return nil
	}

	// Actually start the application

	switch opts.Mode {
	case "cli":

		return process_messages(ctx, opts.MessageIds...)
	case "lambda":

		handler := func(ctx context.Context, sqsEvent events.SQSEvent) error {

			messages := make([]int64, len(sqsEvent.Records))

			for idx, message := range sqsEvent.Records {

				logger := slog.Default()
				logger = logger.With("message id", message.MessageId)

				// logger.Debug("SQS", "message", message.Body)

				var message_id int64

				err := json.Unmarshal([]byte(message.Body), &message_id)

				if err != nil {
					logger.Error("Failed to unmarshal message", "error", err)
					return fmt.Errorf("Failed to unmarshal message ID, %w", err)
				}

				messages[idx] = message_id
			}

			return process_messages(ctx, messages...)
		}

		lambda.Start(handler)

	default:
		return fmt.Errorf("Invalid or unsupported mode")
	}

	return nil
}
