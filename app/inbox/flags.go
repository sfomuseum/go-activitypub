package inbox

import (
	"flag"
	"fmt"
	"os"

	"github.com/sfomuseum/go-flags/flagset"
)

var accounts_database_uri string
var notes_database_uri string
var messages_database_uri string

var account_name string
var verbose bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("follow")

	fs.StringVar(&accounts_database_uri, "accounts-database-uri", "", "A registered sfomuseum/go-activitypub/AccountsDatabase URI.")
	fs.StringVar(&notes_database_uri, "notes-database-uri", "", "A registered sfomuseum/go-activitypub/NotesDatabase URI.")
	fs.StringVar(&messages_database_uri, "messages-database-uri", "", "A registered sfomuseum/go-activitypub/MessagesDatabase URI.")

	fs.StringVar(&account_name, "account-name", "", "The name of the account whose inbox you want to display.")
	fs.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Display all the messges (notes) received for a registered go-activitypub account.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Valid options are:\n")
		fs.PrintDefaults()
	}

	return fs
}
