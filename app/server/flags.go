package server

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

var accounts_database_uri string
var followers_database_uri string
var following_database_uri string
var notes_database_uri string
var messages_database_uri string
var blocks_database_uri string

var server_uri string
var hostname string

var allow_follow bool
var allow_create bool

var verbose bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("activitypub")

	fs.StringVar(&accounts_database_uri, "accounts-database-uri", "", "...")
	fs.StringVar(&followers_database_uri, "followers-database-uri", "", "...")
	fs.StringVar(&following_database_uri, "following-database-uri", "", "...")
	fs.StringVar(&notes_database_uri, "notes-database-uri", "", "...")
	fs.StringVar(&messages_database_uri, "messages-database-uri", "", "...")
	fs.StringVar(&blocks_database_uri, "blocks-database-uri", "", "...")

	fs.BoolVar(&allow_follow, "allow-follow", true, "...")
	fs.BoolVar(&allow_create, "allow-create", false, "...")

	fs.StringVar(&server_uri, "server-uri", "http://localhost:8080", "...")
	fs.StringVar(&hostname, "hostname", "", "...")

	fs.BoolVar(&verbose, "verbose", false, "Enable verbose logging")
	return fs
}
