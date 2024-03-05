package create

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/multi"
)

var accounts_database_uri string
var followers_database_uri string
var posts_database_uri string
var post_tags_database_uri string
var deliveries_database_uri string

var delivery_queue_uri string

var account_name string
var message string
var in_reply_to string

var mentions multi.MultiString

var hostname string
var insecure bool
var verbose bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("follow")

	fs.StringVar(&accounts_database_uri, "accounts-database-uri", "", "...")
	fs.StringVar(&followers_database_uri, "followers-database-uri", "", "...")
	fs.StringVar(&posts_database_uri, "posts-database-uri", "", "...")
	fs.StringVar(&post_tags_database_uri, "post-tags-database-uri", "", "...")
	fs.StringVar(&deliveries_database_uri, "deliveries-database-uri", "", "...")

	fs.StringVar(&delivery_queue_uri, "delivery-queue-uri", "synchronous://", "...")

	fs.StringVar(&account_name, "account-name", "", "...")
	fs.StringVar(&hostname, "hostname", "localhost:8080", "...")
	fs.BoolVar(&insecure, "insecure", false, "...")

	fs.StringVar(&message, "message", "", "...")
	fs.StringVar(&in_reply_to, "in-reply-to", "", "...")
	fs.Var(&mentions, "mention", "...")

	fs.BoolVar(&verbose, "verbose", false, "Enable verbose logging")
	return fs
}
