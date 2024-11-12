package list

import (
	"flag"
	"fmt"
	"os"

	"github.com/sfomuseum/go-flags/flagset"
)

var accounts_database_uri string
var boosts_database_uri string

var account_name string

var hostname string
var insecure bool
var verbose bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("list")

	fs.StringVar(&accounts_database_uri, "accounts-database-uri", "", "A known sfomuseum/go-activitypub/AccountsDatabase URI.")
	fs.StringVar(&boosts_database_uri, "boosts-database-uri", "", "A known sfomuseum/go-activitypub/BlockDatabase URI.")

	fs.StringVar(&account_name, "account-name", "", "The account whose posts have been boosted.")

	fs.StringVar(&hostname, "hostname", "localhost:8080", "The hostname (domain) for the account doing the boosting.")
	fs.BoolVar(&insecure, "insecure", false, "A boolean flag indicating the ActivityPub server delivering activities is insecure (not using TLS).")
	fs.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "List all the boosts received by a go-activitypub account.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Valid options are:\n")
		fs.PrintDefaults()
	}

	return fs
}
