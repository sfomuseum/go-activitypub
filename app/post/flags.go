package post

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

var accounts_database_uri string
var followers_database_uri string
var posts_database_uri string

var delivery_queue_uri string

var account_id string
var hostname string

var message string

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("follow")

	fs.StringVar(&accounts_database_uri, "accounts-database-uri", "", "...")
	fs.StringVar(&followers_database_uri, "followers-database-uri", "", "...")
	fs.StringVar(&posts_database_uri, "posts-database-uri", "", "...")

	fs.StringVar(&delivery_queue_uri, "delivery-queue-uri", "synchronous://", "...")

	fs.StringVar(&account_id, "account-id", "", "...")
	fs.StringVar(&hostname, "hostname", "localhost:8080", "...")

	fs.StringVar(&message, "message", "", "...")

	return fs
}
