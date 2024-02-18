package inbox

import (
	"context"
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

type RunOptions struct {
	AccountsDatabaseURI string
	MessagesDatabaseURI string
	NotesDatabaseURI    string
	AccountName         string
}

func OptionsFromFlagSet(ctx context.Context, fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	opts := &RunOptions{
		AccountsDatabaseURI: accounts_database_uri,
		MessagesDatabaseURI: messages_database_uri,
		NotesDatabaseURI:    notes_database_uri,
		AccountName:         account_name,
	}

	return opts, nil
}
