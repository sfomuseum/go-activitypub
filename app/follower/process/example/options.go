package example

import (
	"context"
	"flag"
	"fmt"

	"github.com/sfomuseum/go-activitypub/uris"
	"github.com/sfomuseum/go-flags/flagset"
)

type RunOptions struct {
	Mode                  string
	FollowerIds           []int64
	MessagesDatabaseURI   string
	NotesDatabaseURI      string
	AccountsDatabaseURI   string
	PropertiesDatabaseURI string
	ActivitiesDatabaseURI string
	PostsDatabaseURI      string
	PostTagsDatabaseURI   string
	DeliveriesDatabaseURI string
	FollowersDatabaseURI  string
	DeliveryQueueURI      string
	MaxAttempts           int
	URIs                  *uris.URIs
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
		Mode:                  mode,
		FollowerIds:           follower_ids,
		MessagesDatabaseURI:   messages_database_uri,
		NotesDatabaseURI:      notes_database_uri,
		AccountsDatabaseURI:   accounts_database_uri,
		PropertiesDatabaseURI: properties_database_uri,
		ActivitiesDatabaseURI: activities_database_uri,
		PostsDatabaseURI:      posts_database_uri,
		PostTagsDatabaseURI:   post_tags_database_uri,
		DeliveriesDatabaseURI: deliveries_database_uri,
		FollowersDatabaseURI:  followers_database_uri,
		DeliveryQueueURI:      delivery_queue_uri,
		URIs:                  uris_table,
		MaxAttempts:           max_attempts,
	}

	return opts, nil
}
