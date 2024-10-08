package deliveries

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

var deliveries_database_uri string
var post_id int64
var recipient string
var verbose bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("follow")

	fs.StringVar(&deliveries_database_uri, "deliveries-database-uri", "", "...")
	fs.Int64Var(&post_id, "post-id", 0, "...")
	fs.StringVar(&recipient, "recipient", "", "...")
	fs.BoolVar(&verbose, "verbose", false, "Enable verbose logging")

	return fs
}
