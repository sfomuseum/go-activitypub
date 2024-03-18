package followers

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

var followers_database_uri string

var verbose bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("activitypub")

	fs.StringVar(&followers_database_uri, "followers-database-uri", "", "...")
	fs.BoolVar(&verbose, "verbose", false, "Enable verbose logging")
	return fs
}
