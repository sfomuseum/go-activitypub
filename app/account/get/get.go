package get

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/sfomuseum/go-activitypub"
	"github.com/sfomuseum/go-activitypub/database"
	"github.com/sfomuseum/go-activitypub/properties"
)

type results struct {
	Account    *activitypub.Account    `json:"account"`
	Properties []*activitypub.Property `json:"properties"`
}

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

	accounts_db, err := database.NewAccountsDatabase(ctx, opts.AccountsDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create accounts database, %w", err)
	}

	defer accounts_db.Close(ctx)

	properties_db, err := database.NewPropertiesDatabase(ctx, opts.PropertiesDatabaseURI)

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

	props_map, err := properties.PropertiesMapForAccount(ctx, properties_db, acct)

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
