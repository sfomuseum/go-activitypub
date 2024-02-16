package server

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

var account_database_uri string
var server_uri string
var hostname string

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("activitypub")

	fs.StringVar(&account_database_uri, "account-database-uri", "mem://accounts/Id", "...")
	fs.StringVar(&server_uri, "server-uri", "http://localhost:8080", "...")
	fs.StringVar(&hostname, "hostname", "", "...")

	return fs
}
