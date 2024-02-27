package add

import (
	"context"
	"flag"
	"fmt"

	"github.com/sfomuseum/go-flags/flagset"
)

type RunOptions struct {
	AccountsDatabaseURI string
	AccountId           int64
	AccountName         string
	AccountType         string
	DisplayName         string
	Blurb               string
	URL                 string
	PublicKeyURI        string
	PrivateKeyURI       string
	AccountIconURI      string
	AllowRemoteIconURI  bool
	EmbedIconURI        bool
}

func OptionsFromFlagSet(ctx context.Context, fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	err := flagset.SetFlagsFromEnvVars(fs, "ACTIVITYPUB")

	if err != nil {
		return nil, fmt.Errorf("Failed to derive flags from environment variables, %w", err)
	}

	opts := &RunOptions{
		AccountsDatabaseURI: accounts_database_uri,
		AccountId:           account_id,
		AccountName:         account_name,
		AccountType:         account_type,
		AccountIconURI:      account_icon_uri,
		AllowRemoteIconURI:  allow_remote_icon_uri,
		EmbedIconURI:        embed_icon_uri,
		DisplayName:         display_name,
		Blurb:               blurb,
		URL:                 account_url,
		PublicKeyURI:        public_key_uri,
		PrivateKeyURI:       private_key_uri,
	}

	return opts, nil
}
