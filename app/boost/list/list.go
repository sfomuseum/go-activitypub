package list

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
		return fmt.Errorf("Failed to create new database, %w", err)
	}

	defer accounts_db.Close(ctx)

	boosts_db, err := database.NewBoostsDatabase(ctx, opts.BoostsDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to instantiate boosts database, %w", err)
	}

	defer boosts_db.Close(ctx)

	logger = logger.With("account", opts.AccountName)

	acct, err := accounts_db.GetAccountWithName(ctx, opts.AccountName)

	if err != nil {
		return fmt.Errorf("Failed to retrieve account %s, %w", opts.AccountName, err)
	}

	cb := func(ctx context.Context, b *activitypub.Boost) error {

		logger.Info("Boost", "id", b.Id, "post", b.PostId, "actor", b.Actor)
		return nil
	}

	err = boosts_db.GetBoostsForAccount(ctx, acct.Id, cb)

	if err != nil {
		return fmt.Errorf("Failed to get boosts, %w", err)
	}

	return nil
}
