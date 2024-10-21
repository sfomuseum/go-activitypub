package list

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

var deliveries_database_uri string
var verbose bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("list")

	fs.StringVar(&deliveries_database_uri, "deliveries-database-uri", "", "...")
	fs.BoolVar(&verbose, "verbose", false, "Enable verbose logging")

	return fs
}
