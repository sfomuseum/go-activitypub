package note

import (
	"context"
	"flag"
	"fmt"

	"github.com/sfomuseum/go-activitypub/uris"
	"github.com/sfomuseum/go-flags/flagset"
)

type RunOptions struct {
	AccountsDatabaseURI   string
	ActivitiesDatabaseURI string
	FollowersDatabaseURI  string
	DeliveriesDatabaseURI string
	DeliveryQueueURI      string
	AccountName           string
	NoteURI               string
	URIs                  *uris.URIs
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
		ActivitiesDatabaseURI: activities_database_uri,
		FollowersDatabaseURI:  followers_database_uri,
		DeliveriesDatabaseURI: deliveries_database_uri,
		DeliveryQueueURI:      delivery_queue_uri,
		AccountName:           account_name,
		NoteURI:               note_uri,
		URIs:                  uris_table,
		Verbose:               verbose,
	}

	return opts, nil
}
