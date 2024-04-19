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
	MaxAttempts           int
	// Allows posts to accounts not followed by author but where account is mentioned in post
	AllowMentions bool
	Verbose       bool
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
		MaxAttempts:           max_attempts,
		Mode:                  mode,
		PostId:                post_id,
		URIs:                  uris_table,
		Verbose:               verbose,
		AllowMentions:         allow_mentions,
	}

	return opts, nil
}
