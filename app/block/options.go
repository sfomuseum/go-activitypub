package block

import (
	"context"
	"flag"
	"fmt"

	"github.com/sfomuseum/go-flags/flagset"
)

type RunOptions struct {
	AccountsDatabaseURI string
	AccountName         string
	BlocksDatabaseURI   string
	BlockName           string
	BlockHost           string
	Undo                bool
}

func OptionsFromFlagSet(ctx context.Context, fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	err := flagset.SetFlagsFromEnvVars(fs, "ACTIVITYPUB")

	if err != nil {
		return nil, fmt.Errorf("Failed to derive flags from environment variables, %w", err)
	}

	opts := &RunOptions{
		AccountsDatabaseURI: accounts_database_uri,
		AccountName:         account_name,
		BlocksDatabaseURI:   blocks_database_uri,
		BlockName:           block_name,
		BlockHost:           block_host,
		Undo:                undo,
	}

	return opts, nil
}
