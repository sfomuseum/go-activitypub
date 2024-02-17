package follow

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

var server_uri string
var hostname string

var accounts_database_uri string
var account_id string

var follow string

// var inbox string

var undo bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("follow")

	fs.StringVar(&server_uri, "server-uri", "http://localhost:8080", "...")
	fs.StringVar(&hostname, "hostname", "", "...")

	fs.StringVar(&accounts_database_uri, "accounts-database-uri", "", "...")
	fs.StringVar(&account_id, "account-id", "", "...")

	fs.StringVar(&follow, "follow", "", "...")
	// fs.StringVar(&inbox, "inbox", "", "...")

	fs.BoolVar(&undo, "undo", false, "...")
	return fs
}
