package post

import (
	"context"
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

type RunOptions struct {
	AccountsDatabaseURI  string
	FollowersDatabaseURI string
	PostsDatabaseURI     string
	DeliveryQueueURI     string
	AccountId            string
}

func OptionsFromFlagSet(ctx context.Context, fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	opts := &RunOptions{
		AccountsDatabaseURI:  accounts_database_uri,
		FollowersDatabaseURI: followers_database_uri,
		PostsDatabaseURI:     posts_database_uri,
		DeliveryQueueURI:     delivery_queue_uri,
		AccountId:            account_id,
	}

	return opts, nil
}
