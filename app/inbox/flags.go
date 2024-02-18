package inbox

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

var accounts_database_uri string
var notes_database_uri string
var messages_database_uri string

var account_name string

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("follow")

	fs.StringVar(&accounts_database_uri, "accounts-database-uri", "", "...")
	fs.StringVar(&notes_database_uri, "notes-database-uri", "", "...")
	fs.StringVar(&messages_database_uri, "messages-database-uri", "", "...")

	fs.StringVar(&account_name, "account-name", "", "...")

	return fs
}
