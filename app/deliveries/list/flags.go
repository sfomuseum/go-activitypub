package list

import (
	"flag"
	"fmt"
	"os"
	
	"github.com/sfomuseum/go-flags/flagset"
)

var deliveries_database_uri string
var verbose bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("list")

	fs.StringVar(&deliveries_database_uri, "deliveries-database-uri", "", "A known sfomuseum/go-activitypub/DeliveriesDatabase URI.")
	fs.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "List all the (ActivityPub activities) deliveries that have been recorded.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Valid options are:\n")
		fs.PrintDefaults()
	}
	
	return fs
}
