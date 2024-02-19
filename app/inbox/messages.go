package inbox

import (
	"context"
	"flag"
	"fmt"
	"log/slog"

	"github.com/sfomuseum/go-activitypub"
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

	slog.SetDefault(logger)

	accounts_db, err := activitypub.NewAccountsDatabase(ctx, opts.AccountsDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to intialize accounts database, %w", err)
	}

	messages_db, err := activitypub.NewMessagesDatabase(ctx, opts.MessagesDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to initialize messages database, %w", err)
	}

	notes_db, err := activitypub.NewNotesDatabase(ctx, opts.NotesDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to initialize notes database, %w", err)
	}

	acct, err := accounts_db.GetAccountWithName(ctx, opts.AccountName)

	if err != nil {
		return fmt.Errorf("Failed to retrieve account %s, %w", opts.AccountName, err)
	}

	messages_cb := func(ctx context.Context, m *activitypub.Message) error {

		logger.Info("Get Note", "message", m.Id, "id", m.NoteId)

		n, err := notes_db.GetNoteWithId(ctx, m.NoteId)

		if err != nil {
			return fmt.Errorf("Failed to retrieve note, %w", err)
		}

		logger.Info("NOTE", "body", string(n.Body))
		return nil
	}

	err = messages_db.GetMessagesForAccount(ctx, acct.Id, messages_cb)

	if err != nil {
		return fmt.Errorf("Failed to retrieve messages, %w", err)
	}

	return nil
}
