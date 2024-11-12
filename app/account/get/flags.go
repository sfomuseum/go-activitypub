package get

import (
	"flag"
	"fmt"
	"os"

	"github.com/sfomuseum/go-flags/flagset"
)

var accounts_database_uri string
var properties_database_uri string

var account_name string

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("activitypub")

	fs.StringVar(&accounts_database_uri, "accounts-database-uri", "null://", "A registered sfomuseum/go-activitypub/AccountsDatabase URI.")
	fs.StringVar(&properties_database_uri, "properties-database-uri", "null://", "A registered sfomuseum/go-activitypub/PropertiesDatabase URI")
	fs.StringVar(&account_name, "account-name", "", "A valid sfomuseum/go-activitypub account name")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Retrieve an ActivityPub account and emit its details as a JSON-encoded string.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Valid options are:\n")
		fs.PrintDefaults()
	}

	return fs
}
