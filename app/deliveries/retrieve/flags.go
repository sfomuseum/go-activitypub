package retrieve

import (
	"flag"
	"fmt"
	"os"
	
	"github.com/sfomuseum/go-flags/flagset"
)

var deliveries_database_uri string
var delivery_id int64

var verbose bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("activitypub")

	fs.StringVar(&deliveries_database_uri, "deliveries-database-uri", "", "A registered sfomuseum/go-activitypub/DeliveriesDatabase URI.")
	fs.Int64Var(&delivery_id, "delivery-id", 0, "The unique ID of the delivery to retrieve.")

	fs.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Retrieve and display a specific (ActivityPub activity) delivery.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Valid options are:\n")
		fs.PrintDefaults()
	}
	
	return fs
}
