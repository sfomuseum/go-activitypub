package add

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/multi"
)

var accounts_database_uri string
var aliases_database_uri string

var account_name string
var aliases multi.MultiCSVString

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("activitypub")

	fs.StringVar(&accounts_database_uri, "accounts-database-uri", "null://", "A known sfomuseum/go-activitypub/AccountsDatabase URI.")
	fs.StringVar(&aliases_database_uri, "aliases-database-uri", "null://", "A known sfomuseum/go-activitypub/AliasesDatabase URI.")
	fs.StringVar(&account_name, "account-name", "", "A valid sfomuseum/go-activitypub account name")
	fs.Var(&aliases, "alias", "One or more aliases to add for an account. Each -alias flag may be a CSV-encoded string containing multiple aliases.")

	return fs
}
