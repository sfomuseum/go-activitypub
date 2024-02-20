package create

import (
	"context"
	"flag"
	"fmt"

	"github.com/sfomuseum/go-flags/flagset"
)

type RunOptions struct {
	Refresh           bool
	DynamodbClientURI string
	TablePrefix       string
}

func OptionsFromFlagSet(ctx context.Context, fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	err := flagset.SetFlagsFromEnvVars(fs, "ACTIVITYPUB")

	if err != nil {
		return nil, fmt.Errorf("Failed to derive flags from environment variables, %w", err)
	}

	opts := &RunOptions{
		Refresh:           refresh,
		DynamodbClientURI: dynamodb_client_uri,
		TablePrefix:       table_prefix,
	}

	return opts, nil
}
