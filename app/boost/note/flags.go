package note

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

var accounts_database_uri string
var followers_database_uri string
var deliveries_database_uri string

var delivery_queue_uri string

var account_name string
var note_uri string

var hostname string
var insecure bool
var verbose bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("follow")

	fs.StringVar(&accounts_database_uri, "accounts-database-uri", "", "...")
	fs.StringVar(&followers_database_uri, "followers-database-uri", "", "...")
	fs.StringVar(&deliveries_database_uri, "deliveries-database-uri", "", "...")

	fs.StringVar(&delivery_queue_uri, "delivery-queue-uri", "synchronous://", "...")

	fs.StringVar(&account_name, "account-name", "", "The account doing the boosting.")
	fs.StringVar(&hostname, "hostname", "localhost:8080", "...")
	fs.BoolVar(&insecure, "insecure", false, "...")

	fs.StringVar(&note_uri, "post", "", "The URI of the note being boosted.")

	fs.BoolVar(&verbose, "verbose", false, "Enable verbose logging")
	return fs
}
