package block

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
		return fmt.Errorf("Failed to initialize accounts database, %w", err)
	}

	defer accounts_db.Close(ctx)

	blocks_db, err := database.NewBlocksDatabase(ctx, opts.BlocksDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to initialize following database, %w", err)
	}

	defer blocks_db.Close(ctx)

	acct, err := accounts_db.GetAccountWithName(ctx, opts.AccountName)

	if err != nil {
		return fmt.Errorf("Failed to retrieve account %s, %w", opts.AccountName, err)
	}

	logger.Info("Process block", "host", opts.BlockHost, "name", opts.BlockName)

	block, err := blocks_db.GetBlockWithAccountIdAndAddress(ctx, acct.Id, opts.BlockHost, opts.BlockName)

	if block != nil {

		if !opts.Undo {
			logger.Info("Block already exists")
			return nil
		}

		err := blocks_db.RemoveBlock(ctx, block)

		if err != nil {
			return fmt.Errorf("Failed to remove block, %w", err)
		}

		logger.Info("Block removed")
		return nil
	}

	block, err = activitypub.NewBlock(ctx, acct.Id, opts.BlockHost, opts.BlockName)

	if err != nil {
		return fmt.Errorf("Failed to create new block, %w", err)
	}

	err = blocks_db.AddBlock(ctx, block)

	if err != nil {
		return fmt.Errorf("Failed to add new block to database, %w", err)
	}

	logger.Info("New block created")
	return nil
}
