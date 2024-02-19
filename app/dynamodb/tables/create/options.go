package create

import (
	"context"
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

type RunOptions struct {
	Refresh           bool
	DynamodbClientURI string
}

func OptionsFromFlagSet(ctx context.Context, fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	opts := &RunOptions{
		Refresh:           refresh,
		DynamodbClientURI: dynamodb_client_uri,
	}

	return opts, nil
}
