package list

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

var activities_database_uri string
var verbose bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("follow")

	fs.StringVar(&activities_database_uri, "activities-database-uri", "", "...")
	fs.BoolVar(&verbose, "verbose", false, "Enable verbose logging")

	return fs
}
