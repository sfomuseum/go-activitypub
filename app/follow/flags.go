package follow

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

var server_uri string
var hostname string

var account_database_uri string
var account_id string

var follow string
var inbox string

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("follow")

	fs.StringVar(&server_uri, "server-uri", "http://localhost:8080", "...")
	fs.StringVar(&hostname, "hostname", "", "...")

	fs.StringVar(&account_database_uri, "account-database-uri", "mem://accounts/Id", "...")
	fs.StringVar(&account_id, "account-id", "", "...")

	fs.StringVar(&follow, "follow", "", "...")
	fs.StringVar(&inbox, "inbox", "", "...")

	return fs
}
