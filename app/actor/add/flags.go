package add

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

var account_database_uri string
var account_id string

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("activitypub")

	fs.StringVar(&account_database_uri, "account-database-uri", "mem://accounts/Id", "...")
	fs.StringVar(&account_id, "account-id", "", "...")

	return fs
}
