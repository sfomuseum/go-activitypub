package list

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

var accounts_database_uri string
var boosts_database_uri string

var account_name string

var hostname string
var insecure bool
var verbose bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("list")

	fs.StringVar(&accounts_database_uri, "accounts-database-uri", "", "...")
	fs.StringVar(&boosts_database_uri, "boosts-database-uri", "", "...")

	fs.StringVar(&account_name, "account-name", "", "The account doing the boosting.")

	fs.StringVar(&hostname, "hostname", "localhost:8080", "...")
	fs.BoolVar(&insecure, "insecure", false, "...")
	fs.BoolVar(&verbose, "verbose", false, "Enable verbose logging")

	return fs
}
