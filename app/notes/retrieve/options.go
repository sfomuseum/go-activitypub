package retrieve

import (
	"context"
	"flag"
	"fmt"

	"github.com/sfomuseum/go-flags/flagset"
)

type RunOptions struct {
	NotesDatabaseURI string
	NoteId           int64
	Body             bool
	Verbose          bool
}

func OptionsFromFlagSet(ctx context.Context, fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	err := flagset.SetFlagsFromEnvVars(fs, "ACTIVITYPUB")

	if err != nil {
		return nil, fmt.Errorf("Failed to derive flags from environment variables, %w", err)
	}

	opts := &RunOptions{
		NotesDatabaseURI: notes_database_uri,
		NoteId:           note_id,
		Body:             body,
		Verbose:          verbose,
	}

	return opts, nil
}
