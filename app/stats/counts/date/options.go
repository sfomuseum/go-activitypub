package date

import (
	"context"
	"flag"
	"fmt"

	"github.com/sfomuseum/go-flags/flagset"
)

type RunOptions struct {
	AccountsDatabaseURI   string
	FollowersDatabaseURI  string
	DeliveriesDatabaseURI string
	FollowingDatabaseURI  string
	NotesDatabaseURI      string
	MessagesDatabaseURI   string
	BlocksDatabaseURI     string
	PostsDatabaseURI      string
	LikesDatabaseURI      string
	BoostsDatabaseURI     string
	Date                  string
	Verbose               bool
}

func OptionsFromFlagSet(ctx context.Context, fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	err := flagset.SetFlagsFromEnvVars(fs, "ACTIVITYPUB")

	if err != nil {
		return nil, fmt.Errorf("Failed to derive flags from environment variables, %w", err)
	}

	opts := &RunOptions{
		AccountsDatabaseURI:   accounts_database_uri,
		DeliveriesDatabaseURI: deliveries_database_uri,
		FollowersDatabaseURI:  followers_database_uri,
		FollowingDatabaseURI:  following_database_uri,
		NotesDatabaseURI:      notes_database_uri,
		MessagesDatabaseURI:   messages_database_uri,
		PostsDatabaseURI:      posts_database_uri,
		BlocksDatabaseURI:     blocks_database_uri,
		LikesDatabaseURI:      likes_database_uri,
		BoostsDatabaseURI:     boosts_database_uri,
		Date:                  date,
		Verbose:               verbose,
	}

	return opts, nil
}
