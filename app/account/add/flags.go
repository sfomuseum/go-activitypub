package add

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

var accounts_database_uri string

var account_name string
var account_id int64

var public_key_uri string
var private_key_uri string

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("activitypub")

	fs.StringVar(&accounts_database_uri, "accounts-database-uri", "", "...")

	fs.StringVar(&account_name, "account-name", "", "...")
	fs.Int64Var(&account_id, "account-id", 0, "...")

	fs.StringVar(&public_key_uri, "public-key-uri", "", "...")
	fs.StringVar(&private_key_uri, "private-key-uri", "", "...")
	return fs
}
