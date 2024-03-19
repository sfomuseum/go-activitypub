package deliveries

import (
	"context"
	"flag"
	"fmt"

	"github.com/sfomuseum/go-flags/flagset"
)

type RunOptions struct {
	DeliveriesDatabaseURI string
	Recipient             string
	PostId                int64
	Verbose               bool
}

func OptionsFromFlagSet(ctx context.Context, fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	err := flagset.SetFlagsFromEnvVars(fs, "ACTIVITYPUB")

	if err != nil {
		return nil, fmt.Errorf("Failed to derive flags from environment variables, %w", err)
	}

	opts := &RunOptions{
		DeliveriesDatabaseURI: deliveries_database_uri,
		Recipient:             recipient,
		PostId:                post_id,
		Verbose:               verbose,
	}

	return opts, nil
}
