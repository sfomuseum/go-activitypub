package deliver

import (
	"context"
	"flag"
	"fmt"

	"github.com/sfomuseum/go-activitypub/uris"
	"github.com/sfomuseum/go-flags/flagset"
)

type RunOptions struct {
	AccountsDatabaseURI   string
	FollowersDatabaseURI  string
	PostsDatabaseURI      string
	PostTagsDatabaseURI   string
	DeliveriesDatabaseURI string
	DeliveryQueueURI      string
	URIs                  *uris.URIs
	Mode                  string
	PostId                int64
	Verbose               bool
}

func OptionsFromFlagSet(ctx context.Context, fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	err := flagset.SetFlagsFromEnvVars(fs, "ACTIVITYPUB")

	if err != nil {
		return nil, fmt.Errorf("Failed to derive flags from environment variables, %w", err)
	}

	uris_table := uris.DefaultURIs()
	uris_table.Hostname = hostname
	uris_table.Insecure = insecure

	opts := &RunOptions{
		AccountsDatabaseURI:   accounts_database_uri,
		FollowersDatabaseURI:  followers_database_uri,
		PostsDatabaseURI:      posts_database_uri,
		PostTagsDatabaseURI:   post_tags_database_uri,
		DeliveriesDatabaseURI: deliveries_database_uri,
		DeliveryQueueURI:      delivery_queue_uri,
		Mode:                  mode,
		PostId:                post_id,
		URIs:                  uris_table,
		Verbose:               verbose,
	}

	return opts, nil
}
