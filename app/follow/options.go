package follow

import (
	"context"
	"flag"
	"fmt"

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
	Verbose              bool
}

func OptionsFromFlagSet(ctx context.Context, fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	err := flagset.SetFlagsFromEnvVars(fs, "ACTIVITYPUB")

	if err != nil {
		return nil, fmt.Errorf("Failed to derive flags from environment variables, %w", err)
	}

	uris_table := activitypub.DefaultURIs()
	uris_table.Hostname = hostname
	uris_table.Insecure = insecure

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
		Verbose:              verbose,
	}

	return opts, nil
}
