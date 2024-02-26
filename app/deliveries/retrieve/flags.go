package retrieve

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

var deliveries_database_uri string
var delivery_id int64

var verbose bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("activitypub")

	fs.StringVar(&deliveries_database_uri, "deliveries-database-uri", "", "...")
	fs.Int64Var(&delivery_id, "delivery-id", 0, "...")

	fs.BoolVar(&verbose, "verbose", false, "Enable verbose logging")
	return fs
}
