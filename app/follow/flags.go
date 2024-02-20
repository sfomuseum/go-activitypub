package follow

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

var server_uri string
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

	fs.StringVar(&server_uri, "server-uri", "http://localhost:8080", "...")
	fs.StringVar(&hostname, "hostname", "localhost:8080", "...")
	fs.BoolVar(&insecure, "insecure", false, "...")

	fs.StringVar(&accounts_database_uri, "accounts-database-uri", "", "...")
	fs.StringVar(&following_database_uri, "following-database-uri", "", "...")
	fs.StringVar(&messages_database_uri, "messages-database-uri", "", "...")

	fs.StringVar(&account_name, "account-name", "", "...")

	fs.StringVar(&follow_address, "follow", "", "...")

	fs.BoolVar(&undo, "undo", false, "...")

	fs.BoolVar(&verbose, "verbose", false, "Enable verbose logging")
	return fs
}
