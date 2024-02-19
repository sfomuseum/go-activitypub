package follow

import (
	"context"
	"flag"

	"github.com/sfomuseum/go-activitypub"
	"github.com/sfomuseum/go-flags/flagset"
)

type RunOptions struct {
	ServerURI            string
	Hostname             string
	AccountsDatabaseURI  string
	FollowingDatabaseURI string
	MessagesDatabaseURI  string
	AccountName          string
	FollowAddress        string
	Undo                 bool
	URIs                 *activitypub.URIs
}

func OptionsFromFlagSet(ctx context.Context, fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	uris_table := activitypub.DefaultURIs()

	opts := &RunOptions{
		ServerURI:            server_uri,
		Hostname:             hostname,
		AccountsDatabaseURI:  accounts_database_uri,
		FollowingDatabaseURI: following_database_uri,
		MessagesDatabaseURI:  messages_database_uri,
		AccountName:          account_name,
		FollowAddress:        follow_address,
		Undo:                 undo,
		URIs:                 uris_table,
	}

	return opts, nil
}
