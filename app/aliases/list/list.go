// Gather the list of aliases for a given account and emit as JSON-encoded string to STDOUT.
package list

/*

$> go run cmd/list-aliases/main.go \
	-accounts-database-uri 'awsdynamodb://accounts?partition_key=Id&allow_scans=true&region=us-west-2&credentials=session' \
	-aliases-database-uri 'awsdynamodb://aliases?partition_key=Name&allow_scans=true&region=us-west-2&credentials=session' \
	-account-name 102527513 \
| jq

[
  {
    "name": "SFOairport",
    "account_id": 1765465178285019136,
    "created": 1709754673
  },
  {
    "name": "KSFOairport",
    "account_id": 1765465178285019136,
    "created": 1709754674
  }
]

*/

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/sfomuseum/go-activitypub"
	"github.com/sfomuseum/go-activitypub/database"
	"github.com/sfomuseum/go-activitypub/slog"
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

	logger := slog.Default()

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
