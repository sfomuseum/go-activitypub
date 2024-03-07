package date

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

var accounts_database_uri string
var blocks_database_uri string
var boosts_database_uri string
var followers_database_uri string
var following_database_uri string
var likes_database_uri string
var notes_database_uri string
var messages_database_uri string
var posts_database_uri string
var deliveries_database_uri string

var date string
var verbose bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("activitypub")

	fs.StringVar(&accounts_database_uri, "accounts-database-uri", "", "...")
	fs.StringVar(&followers_database_uri, "followers-database-uri", "", "...")
	fs.StringVar(&following_database_uri, "following-database-uri", "", "...")
	fs.StringVar(&notes_database_uri, "notes-database-uri", "", "...")
	fs.StringVar(&messages_database_uri, "messages-database-uri", "", "...")
	fs.StringVar(&blocks_database_uri, "blocks-database-uri", "", "...")
	fs.StringVar(&posts_database_uri, "posts-database-uri", "", "...")
	fs.StringVar(&likes_database_uri, "likes-database-uri", "", "...")
	fs.StringVar(&boosts_database_uri, "boosts-database-uri", "", "...")
	fs.StringVar(&deliveries_database_uri, "deliveries-database-uri", "", "...")

	fs.StringVar(&date, "date", "", "...")
	fs.BoolVar(&verbose, "verbose", false, "Enable verbose logging")
	return fs
}
