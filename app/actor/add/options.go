package add

import (
	"context"
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

type RunOptions struct {
	AccountDatabaseURI string
	AccountId          string
}

func OptionsFromFlagSet(ctx context.Context, fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	opts := &RunOptions{
		AccountDatabaseURI: account_database_uri,
		AccountId:          account_id,
	}

	return opts, nil
}
