package list

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

var accounts_database_uri string
var aliases_database_uri string

var account_name string

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("activitypub")

	fs.StringVar(&accounts_database_uri, "accounts-database-uri", "null://", "A known sfomuseum/go-activitypub/AccountsDatabase URI.")
	fs.StringVar(&aliases_database_uri, "aliases-database-uri", "null://", "A known sfomuseum/go-activitypub/AliasesDatabase URI.")

	fs.StringVar(&account_name, "account-name", "", "A valid sfomuseum/go-activitypub account name")
	return fs
}
