package block

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
		return fmt.Errorf("Failed to initialize accounts database, %w", err)
	}

	blocks_db, err := activitypub.NewBlocksDatabase(ctx, opts.BlocksDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to initialize following database, %w", err)
	}

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
