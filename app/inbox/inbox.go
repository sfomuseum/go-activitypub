package inbox

import (
	"context"
	"flag"
	"fmt"
	"log/slog"

	"github.com/sfomuseum/go-activitypub"
	"github.com/sfomuseum/go-activitypub/database"
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
		return fmt.Errorf("Failed to intialize accounts database, %w", err)
	}

	defer accounts_db.Close(ctx)

	messages_db, err := database.NewMessagesDatabase(ctx, opts.MessagesDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to initialize messages database, %w", err)
	}

	defer messages_db.Close(ctx)

	notes_db, err := database.NewNotesDatabase(ctx, opts.NotesDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to initialize notes database, %w", err)
	}

	defer notes_db.Close(ctx)

	logger = logger.With("account", opts.AccountName)

	acct, err := accounts_db.GetAccountWithName(ctx, opts.AccountName)

	if err != nil {
		return fmt.Errorf("Failed to retrieve account %s, %w", opts.AccountName, err)
	}

	logger = logger.With("account id", acct.Id)

	messages_cb := func(ctx context.Context, m *activitypub.Message) error {

		logger.Info("Get Note", "message", m.Id, "id", m.NoteId)

		n, err := notes_db.GetNoteWithId(ctx, m.NoteId)

		if err != nil {
			return fmt.Errorf("Failed to retrieve note, %w", err)
		}

		logger.Info("NOTE", "body", string(n.Body))
		return nil
	}

	logger.Debug("Get messages")

	err = messages_db.GetMessagesForAccount(ctx, acct.Id, messages_cb)

	if err != nil {
		return fmt.Errorf("Failed to retrieve messages, %w", err)
	}

	return nil
}
