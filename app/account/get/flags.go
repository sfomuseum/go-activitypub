package get

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

var accounts_database_uri string
var properties_database_uri string

var account_name string

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("activitypub")

	fs.StringVar(&accounts_database_uri, "accounts-database-uri", "", "...")
	fs.StringVar(&properties_database_uri, "properties-database-uri", "", "...")

	fs.StringVar(&account_name, "account-name", "", "...")
	return fs
}
