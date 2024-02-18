package post

import (
	"context"
	"flag"

	"github.com/sfomuseum/go-activitypub"
	"github.com/sfomuseum/go-flags/flagset"
)

type RunOptions struct {
	AccountsDatabaseURI  string
	FollowersDatabaseURI string
	PostsDatabaseURI     string
	DeliveryQueueURI     string
	AccountId            string
	Message              string
	Hostname             string
	URIs                 *activitypub.URIs
}

func OptionsFromFlagSet(ctx context.Context, fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	uris_table := activitypub.DefaultURIs()

	opts := &RunOptions{
		AccountsDatabaseURI:  accounts_database_uri,
		FollowersDatabaseURI: followers_database_uri,
		PostsDatabaseURI:     posts_database_uri,
		DeliveryQueueURI:     delivery_queue_uri,
		AccountId:            account_id,
		Hostname:             hostname,
		Message:              message,
		URIs:                 uris_table,
	}

	return opts, nil
}
