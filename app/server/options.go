package server

import (
	"context"
	"flag"
	"fmt"
	"html/template"
	"net/url"

	"github.com/mitchellh/copystructure"
	"github.com/sfomuseum/go-activitypub/templates/html"
	"github.com/sfomuseum/go-activitypub/uris"
	"github.com/sfomuseum/go-flags/flagset"
)

type RunOptions struct {
	ServerURI            string
	URIs                 *uris.URIs
	AccountsDatabaseURI  string
	AliasesDatabaseURI   string
	FollowersDatabaseURI string
	FollowingDatabaseURI string
	NotesDatabaseURI     string
	MessagesDatabaseURI  string
	BlocksDatabaseURI    string
	PostsDatabaseURI     string
	LikesDatabaseURI     string
	BoostsDatabaseURI    string
	RepliesDatabaseURI   string
	AllowFollow          bool
	AllowCreate          bool
	AllowLikes           bool
	AllowBoosts          bool
	AllowReplies         bool
	AllowRemoteIconURI   bool
	Verbose              bool
	Templates            *template.Template
}

func OptionsFromFlagSet(ctx context.Context, fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	err := flagset.SetFlagsFromEnvVars(fs, "ACTIVITYPUB")

	if err != nil {
		return nil, fmt.Errorf("Failed to derive flags from environment variables, %w", err)
	}

	if hostname == "" {

		u, err := url.Parse(server_uri)

		if err != nil {
			return nil, fmt.Errorf("Failed to parse server URI, %w", err)
		}

		hostname = u.Host
	}

	uris_table := uris.DefaultURIs()
	uris_table.Hostname = hostname
	uris_table.Insecure = insecure

	t, err := html.LoadTemplates(ctx)

	if err != nil {
		return nil, fmt.Errorf("Failed to load templates, %w", err)
	}

	opts := &RunOptions{
		AccountsDatabaseURI:  accounts_database_uri,
		AliasesDatabaseURI:   aliases_database_uri,
		FollowersDatabaseURI: followers_database_uri,
		FollowingDatabaseURI: following_database_uri,
		NotesDatabaseURI:     notes_database_uri,
		MessagesDatabaseURI:  messages_database_uri,
		PostsDatabaseURI:     posts_database_uri,
		BlocksDatabaseURI:    blocks_database_uri,
		LikesDatabaseURI:     likes_database_uri,
		BoostsDatabaseURI:    boosts_database_uri,
		RepliesDatabaseURI:   replies_database_uri,
		ServerURI:            server_uri,
		URIs:                 uris_table,
		AllowFollow:          allow_follow,
		AllowCreate:          allow_create,
		AllowBoosts:          allow_boosts,
		AllowLikes:           allow_likes,
		AllowReplies:         allow_replies,
		AllowRemoteIconURI:   allow_remote_icon_uri,
		Verbose:              verbose,
		Templates:            t,
	}

	return opts, nil
}

// Is this (clone) really necessary? I am starting to think it is not...

func (o *RunOptions) clone() (*RunOptions, error) {

	v, err := copystructure.Copy(o)

	if err != nil {
		return nil, fmt.Errorf("Failed to create local run options, %w", err)
	}

	new_opts := v.(*RunOptions)

	new_opts.Templates = o.Templates
	return new_opts, nil
}
