package list

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/sfomuseum/go-activitypub"
)

type results struct {
	Account    *activitypub.Account    `json:"account"`
	Properties []*activitypub.Property `json:"properties"`
}

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
		return fmt.Errorf("Failed to create accounts database, %w", err)
	}

	defer accounts_db.Close(ctx)

	aliases_db, err := activitypub.NewAliasesDatabase(ctx, opts.AliasesDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create aliases database, %w", err)
	}

	defer aliases_db.Close(ctx)

	acct, err := accounts_db.GetAccountWithName(ctx, opts.AccountName)

	if err != nil {
		return fmt.Errorf("Failed to retrieve account %s, %w", opts.AccountName, err)
	}

	aliases := make([]*activitypub.Alias, 0)

	aliases_cb := func(ctx context.Context, a *activitypub.Alias) error {
		aliases = append(aliases, a)
		return nil
	}

	err = aliases_db.GetAliasesForAccount(ctx, acct.Id, aliases_cb)

	if err != nil {
		return fmt.Errorf("Failed to retrieve aliases for account, %w", err)
	}

	enc := json.NewEncoder(os.Stdout)
	err = enc.Encode(aliases)

	if err != nil {
		return fmt.Errorf("Failed to encode results, %w", err)
	}

	return nil
}
