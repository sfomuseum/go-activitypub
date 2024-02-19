package create

import (
	"context"
	"flag"
)

type RunOptions struct {
	Refresh           bool
	DynamodbClientURI string
}

func OptionsFromFlagSet(ctx context.Context, fs *flag.FlagSet) (*RunOptions, error) {

	opts := &RunOptions{
		Refresh:           refresh,
		DynamodbClientURI: dynamodb_client_uri,
	}

	return opts, nil
}
