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
	ServerURI            string
	AccountsDatabaseURI  string
	FollowersDatabaseURI string
	FollowingDatabaseURI string
	NotesDatabaseURI     string
	MessagesDatabaseURI  string
	BlocksDatabaseURI    string
	URIs                 *activitypub.URIs
	AllowFollow          bool
	AllowCreate          bool
	Verbose              bool
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
	uris_table.Hostname = hostname
	uris_table.Insecure = insecure

	opts := &RunOptions{
		AccountsDatabaseURI:  accounts_database_uri,
		FollowersDatabaseURI: followers_database_uri,
		FollowingDatabaseURI: following_database_uri,
		NotesDatabaseURI:     notes_database_uri,
		MessagesDatabaseURI:  messages_database_uri,
		BlocksDatabaseURI:    blocks_database_uri,
		ServerURI:            server_uri,
		URIs:                 uris_table,
		AllowFollow:          allow_follow,
		AllowCreate:          allow_create,
		Verbose:              verbose,
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
