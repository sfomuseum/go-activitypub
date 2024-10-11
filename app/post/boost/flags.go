package boost

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

var accounts_database_uri string
var followers_database_uri string

// var posts_database_uri string
// var post_tags_database_uri string
var deliveries_database_uri string

var delivery_queue_uri string

var account_name string
var post string

// var max_attempts int

var hostname string
var insecure bool
var verbose bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("follow")

	fs.StringVar(&accounts_database_uri, "accounts-database-uri", "", "...")
	fs.StringVar(&followers_database_uri, "followers-database-uri", "", "...")
	// fs.StringVar(&posts_database_uri, "posts-database-uri", "", "...")
	// fs.StringVar(&post_tags_database_uri, "post-tags-database-uri", "", "...")
	fs.StringVar(&deliveries_database_uri, "deliveries-database-uri", "", "...")

	fs.StringVar(&delivery_queue_uri, "delivery-queue-uri", "synchronous://", "...")

	fs.StringVar(&account_name, "account-name", "", "The account doing the boosting.")
	fs.StringVar(&hostname, "hostname", "localhost:8080", "...")
	fs.BoolVar(&insecure, "insecure", false, "...")

	// fs.IntVar(&max_attempts, "max-attempts", 5, "...")

	fs.StringVar(&post, "post", "", "The URI of the post being boosted.")

	fs.BoolVar(&verbose, "verbose", false, "Enable verbose logging")
	return fs
}
