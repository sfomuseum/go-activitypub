package follow

import (
	"context"
	"flag"

	"github.com/sfomuseum/go-activitypub"
	"github.com/sfomuseum/go-flags/flagset"
)

type RunOptions struct {
	ServerURI          string
	Hostname           string
	AccountDatabaseURI string
	AccountId          string
	Follow             string
	Inbox              string
	URIs               *activitypub.URIs
}

func OptionsFromFlagSet(ctx context.Context, fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	uris_table := activitypub.DefaultURIs()

	opts := &RunOptions{
		ServerURI:          server_uri,
		Hostname:           hostname,
		AccountDatabaseURI: account_database_uri,
		AccountId:          account_id,
		Follow:             follow,
		Inbox:              inbox,
		URIs:               uris_table,
	}

	return opts, nil
}
