package example

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/multi"
)

var mode string
var message_ids multi.MultiInt64

var accounts_database_uri string
var properties_database_uri string
var activities_database_uri string
var messages_database_uri string
var notes_database_uri string
var posts_database_uri string
var post_tags_database_uri string
var deliveries_database_uri string
var followers_database_uri string

var delivery_queue_uri string

var max_attempts int
var hostname string
var insecure bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("activitypub")

	fs.StringVar(&mode, "mode", "cli", "Valid options are: cli, lambda")
	fs.Var(&message_ids, "id", "...")

	fs.StringVar(&accounts_database_uri, "accounts-database-uri", "", "...")
	fs.StringVar(&properties_database_uri, "properties-database-uri", "", "...")
	fs.StringVar(&activities_database_uri, "activities-database-uri", "", "...")
	fs.StringVar(&messages_database_uri, "messages-database-uri", "", "...")
	fs.StringVar(&notes_database_uri, "notes-database-uri", "", "...")
	fs.StringVar(&posts_database_uri, "posts-database-uri", "", "...")
	fs.StringVar(&post_tags_database_uri, "post-tags-database-uri", "", "...")
	fs.StringVar(&deliveries_database_uri, "deliveries-database-uri", "", "...")
	fs.StringVar(&followers_database_uri, "followers-database-uri", "", "...")

	fs.StringVar(&delivery_queue_uri, "delivery-queue-uri", "synchronous://", "...")

	fs.IntVar(&max_attempts, "max-attempts", 5, "...")
	fs.StringVar(&hostname, "hostname", "", "...")
	fs.BoolVar(&insecure, "insecure", false, "...")

	return fs
}
