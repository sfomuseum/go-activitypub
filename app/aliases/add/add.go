// Add one or more aliases for a sfomuseum/go-activitypub account.
package add

/*

$> go run cmd/add-aliases/main.go \
	-accounts-database-uri 'awsdynamodb://accounts?partition_key=Id&allow_scans=true&region=us-west-2&credentials=session' \
	-aliases-database-uri 'awsdynamodb://aliases?partition_key=Name&allow_scans=true&region=us-west-2&credentials=session' \
	-account-name 1762688757 \
	-alias 2015.166.1180

time=2024-10-04T10:35:36.480-07:00 level=INFO msg="New alias created" account=1762688757 "account ID"=1777551842624933888 alias=2015.166.1180

*/

import (
	"context"
	"flag"
	"fmt"
	"slices"
	"time"

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

	accounts_db, err := database.NewAccountsDatabase(ctx, opts.AccountsDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create accounts database, %w", err)
	}

	defer accounts_db.Close(ctx)

	aliases_db, err := database.NewAliasesDatabase(ctx, opts.AliasesDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create aliases database, %w", err)
	}

	defer aliases_db.Close(ctx)

	acct, err := accounts_db.GetAccountWithName(ctx, opts.AccountName)

	if err != nil {
		return fmt.Errorf("Failed to retrieve account %s, %w", opts.AccountName, err)
	}

	if len(opts.Aliases) == 0 {
		logger.Debug("No aliases to add")
		return nil
	}

	logger = logger.With("account", opts.AccountName)
	logger = logger.With("account ID", acct.Id)

	current_aliases := make([]string, 0)

	aliases_cb := func(ctx context.Context, a *activitypub.Alias) error {
		current_aliases = append(current_aliases, a.Name)
		return nil
	}

	err = aliases_db.GetAliasesForAccount(ctx, acct.Id, aliases_cb)

	if err != nil {
		return fmt.Errorf("Failed to retrieve aliases for account, %w", err)
	}

	for _, a := range opts.Aliases {

		if slices.Contains(current_aliases, a) {
			logger.Info("Alias already registered for account, skipping", "alias", a)
			continue
		}

		taken, err := aliases.IsAliasNameTaken(ctx, aliases_db, a)

		if err != nil {
			return fmt.Errorf("Failed to determine if alias (%s) is taken, %w", a, err)
		}

		if taken {
			return fmt.Errorf("Alias (%s) is already taken", a)
		}

		now := time.Now()
		ts := now.Unix()

		new_alias := &activitypub.Alias{
			Name:      a,
			AccountId: acct.Id,
			Created:   ts,
		}

		err = aliases_db.AddAlias(ctx, new_alias)

		if err != nil {
			return fmt.Errorf("Failed to add new alias for %s, %w", a, err)
		}

		logger.Info("New alias created", "alias", new_alias.Name)
	}

	return nil
}
