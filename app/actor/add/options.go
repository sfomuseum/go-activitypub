package add

import (
	"context"
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

type RunOptions struct {
	DatabaseURI string
	ActorId     string
}

func OptionsFromFlagSet(ctx context.Context, fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	opts := &RunOptions{
		DatabaseURI: database_uri,
		ActorId:     actor_id,
	}

	return opts, nil
}
