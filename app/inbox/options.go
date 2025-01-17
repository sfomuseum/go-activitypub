package inbox

import (
	"context"
	"flag"
	"fmt"

	"github.com/sfomuseum/go-flags/flagset"
)

type RunOptions struct {
	AccountsDatabaseURI string
	MessagesDatabaseURI string
	NotesDatabaseURI    string
	AccountName         string
	Verbose             bool
}

func OptionsFromFlagSet(ctx context.Context, fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	err := flagset.SetFlagsFromEnvVars(fs, "ACTIVITYPUB")

	if err != nil {
		return nil, fmt.Errorf("Failed to derive flags from environment variables, %w", err)
	}

	opts := &RunOptions{
		AccountsDatabaseURI: accounts_database_uri,
		MessagesDatabaseURI: messages_database_uri,
		NotesDatabaseURI:    notes_database_uri,
		AccountName:         account_name,
		Verbose:             verbose,
	}

	return opts, nil
}
