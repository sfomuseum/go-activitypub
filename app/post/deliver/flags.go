package deliver

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

var accounts_database_uri string
var followers_database_uri string
var posts_database_uri string
var post_tags_database_uri string
var deliveries_database_uri string

var delivery_queue_uri string
var max_attempts int

var post_id int64
var mode string

var hostname string
var insecure bool
var verbose bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("follow")

	fs.StringVar(&accounts_database_uri, "accounts-database-uri", "", "...")
	fs.StringVar(&followers_database_uri, "followers-database-uri", "", "...")
	fs.StringVar(&posts_database_uri, "posts-database-uri", "", "...")
	fs.StringVar(&post_tags_database_uri, "post-tags-database-uri", "null://", "...")
	fs.StringVar(&deliveries_database_uri, "deliveries-database-uri", "", "...")

	fs.StringVar(&delivery_queue_uri, "delivery-queue-uri", "synchronous://", "...")

	fs.IntVar(&max_attempts, "max-attempts", 0, "...")
	fs.Int64Var(&post_id, "post-id", 0, "...")

	fs.StringVar(&mode, "mode", "cli", "...")

	fs.StringVar(&hostname, "hostname", "localhost:8080", "...")
	fs.BoolVar(&insecure, "insecure", false, "...")
	fs.BoolVar(&verbose, "verbose", false, "Enable verbose logging")

	return fs
}
