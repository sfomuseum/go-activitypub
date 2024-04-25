package get

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

	properties_db, err := activitypub.NewPropertiesDatabase(ctx, opts.PropertiesDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create properties database, %w", err)
	}

	defer properties_db.Close(ctx)

	r := new(results)

	acct, err := accounts_db.GetAccountWithName(ctx, opts.AccountName)

	if err != nil {
		return fmt.Errorf("Failed to retrieve account %s, %w", opts.AccountName, err)
	}

	r.Account = acct

	props_map, err := activitypub.PropertiesMapForAccount(ctx, properties_db, acct)

	if err != nil {
		return fmt.Errorf("Failed to derive properties for account %s, %w", opts.AccountName, err)
	}

	props := make([]*activitypub.Property, 0)

	for _, pr := range props_map {
		props = append(props, pr)
	}

	r.Properties = props

	enc := json.NewEncoder(os.Stdout)
	err = enc.Encode(r)

	if err != nil {
		return fmt.Errorf("Failed to encode results, %w", err)
	}

	return nil
}
