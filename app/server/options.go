package server

import (
	"context"
	"flag"
	"fmt"

	"github.com/mitchellh/copystructure"
	"github.com/sfomuseum/go-flags/flagset"
)

type RunOptions struct {
	ServerURI        string
	ActorDatabaseURI string
}

func OptionsFromFlagSet(ctx context.Context, fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	opts := &RunOptions{
		ActorDatabaseURI: actor_database_uri,
		ServerURI:        server_uri,
	}

	return opts, nil
}

func (o *RunOptions) clone() (*RunOptions, error) {

	v, err := copystructure.Copy(o)

	if err != nil {
		return nil, fmt.Errorf("Failed to create local run options, %w", err)
	}

	new_opts := v.(*RunOptions)
	return new_opts, nil
}
