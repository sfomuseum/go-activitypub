package server

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

var accounts_database_uri string
var followers_database_uri string
var following_database_uri string

var server_uri string
var hostname string

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("activitypub")

	fs.StringVar(&accounts_database_uri, "accounts-database-uri", "", "...")
	fs.StringVar(&followers_database_uri, "followers-database-uri", "", "...")
	fs.StringVar(&following_database_uri, "following-database-uri", "", "...")

	fs.StringVar(&server_uri, "server-uri", "http://localhost:8080", "...")
	fs.StringVar(&hostname, "hostname", "", "...")

	return fs
}
