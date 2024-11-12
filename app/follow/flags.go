package follow

import (
	"flag"
	"fmt"
	"os"

	"github.com/sfomuseum/go-flags/flagset"
)

var hostname string
var insecure bool

var accounts_database_uri string
var following_database_uri string
var messages_database_uri string

var account_name string
var follow_address string

var undo bool
var verbose bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("follow")

	fs.StringVar(&accounts_database_uri, "accounts-database-uri", "", "A registered sfomuseum/go-activitypub/AccountsDatabase URI.")
	fs.StringVar(&following_database_uri, "following-database-uri", "", "A registered sfomuseum/go-activitypub/FollowingDatabase URI.")
	fs.StringVar(&messages_database_uri, "messages-database-uri", "", "A registered sfomuseum/go-activitypub/MessagesDatabase URI.")

	fs.StringVar(&account_name, "account-name", "", "The name of the account doing the following.")

	fs.StringVar(&follow_address, "follow", "", "The ActivityPub @user@host address to follow.")

	fs.BoolVar(&undo, "undo", false, "Stop following the account defined by the -follow flag.")

	fs.StringVar(&hostname, "hostname", "localhost:8080", "The hostname (domain) of the ActivityPub server delivering activities.")
	fs.BoolVar(&insecure, "insecure", false, "A boolean flag indicating the ActivityPub server delivering activities is insecure (not using TLS).")

	fs.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Follow a @user@host ActivityPub account on behalf of a registered go-activitypub account.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Valid options are:\n")
		fs.PrintDefaults()
	}

	return fs
}
