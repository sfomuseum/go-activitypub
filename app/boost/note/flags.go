package note

import (
	"flag"
	"fmt"
	"os"

	"github.com/sfomuseum/go-flags/flagset"
)

var accounts_database_uri string
var activities_database_uri string
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

	fs.StringVar(&accounts_database_uri, "accounts-database-uri", "", "A known sfomuseum/go-activitypub/AccountsDatabase URI.")
	fs.StringVar(&activities_database_uri, "activities-database-uri", "", "A known sfomuseum/go-activitypub/ActivitiesDatabase URI.")
	fs.StringVar(&followers_database_uri, "followers-database-uri", "", "A known sfomuseum/go-activitypub/FollowersDatabase URI.")
	fs.StringVar(&deliveries_database_uri, "deliveries-database-uri", "", "A known sfomuseum/go-activitypub/DeliveriesDatabase URI.")

	fs.StringVar(&delivery_queue_uri, "delivery-queue-uri", "synchronous://", "A known sfomuseum/go-activitypub/queue.DeliveryQueue URI.")

	fs.StringVar(&account_name, "account-name", "", "The account doing the boosting.")
	fs.StringVar(&hostname, "hostname", "localhost:8080", "The hostname (domain) for the account doing the boosting.")
	fs.BoolVar(&insecure, "insecure", false, "A boolean flag indicating the ActivityPub server delivering activities is insecure (not using TLS).")
	fs.StringVar(&note_uri, "note", "", "The URI of the note being boosted.")
	fs.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Boost an ActivityPub note on behalf of a registered go-activity account.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Valid options are:\n")
		fs.PrintDefaults()
	}
	
	return fs
}
