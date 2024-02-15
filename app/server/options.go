package server

import (
	"context"
	"flag"
	"fmt"
	"net/url"

	"github.com/mitchellh/copystructure"
	"github.com/sfomuseum/go-activitypub"
	"github.com/sfomuseum/go-flags/flagset"
)

type RunOptions struct {
	ServerURI        string
	ActorDatabaseURI string
	Hostname         string
	URIs             *activitypub.URIs
}

func OptionsFromFlagSet(ctx context.Context, fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	if hostname == "" {

		u, err := url.Parse(server_uri)

		if err != nil {
			return nil, fmt.Errorf("Failed to parse server URI, %w", err)
		}

		hostname = u.Host
	}

	uris_table := activitypub.DefaultURIs()

	opts := &RunOptions{
		ActorDatabaseURI: actor_database_uri,
		ServerURI:        server_uri,
		Hostname:         hostname,
		URIs:             uris_table,
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
