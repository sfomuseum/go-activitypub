package add

import (
	"context"
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

type RunOptions struct {
	AccountsDatabaseURI string
	AccountName         string
	AccountId           int64
	PublicKeyURI        string
	PrivateKeyURI       string
}

func OptionsFromFlagSet(ctx context.Context, fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	opts := &RunOptions{
		AccountsDatabaseURI: accounts_database_uri,
		AccountId:           account_id,
		AccountName:         account_name,
		PublicKeyURI:        public_key_uri,
		PrivateKeyURI:       private_key_uri,
	}

	return opts, nil
}
