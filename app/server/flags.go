package server

import (
	"flag"
	"fmt"
	"os"

	"github.com/sfomuseum/go-flags/flagset"
)

var accounts_database_uri string
var aliases_database_uri string
var followers_database_uri string
var following_database_uri string
var notes_database_uri string
var messages_database_uri string
var posts_database_uri string
var post_tags_database_uri string
var properties_database_uri string
var blocks_database_uri string
var likes_database_uri string
var boosts_database_uri string

var process_message_queue_uri string
var process_follow_queue_uri string

var server_uri string
var hostname string
var insecure bool

var allow_follow bool
var allow_create bool
var allow_likes bool
var allow_boosts bool

// Allows posts to accounts not followed by author but where account is mentioned in post
var allow_mentions bool

var allow_remote_icon_uri bool

var disabled bool
var verbose bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("activitypub")

	fs.StringVar(&accounts_database_uri, "accounts-database-uri", "", "A registered sfomuseum/go-activitypub/database.DeliveriesDatabase URI.")
	fs.StringVar(&aliases_database_uri, "aliases-database-uri", "", "A registered sfomuseum/go-activitypub/database.AliasesDatabase URI.")
	fs.StringVar(&followers_database_uri, "followers-database-uri", "", "A registered sfomuseum/go-activitypub/database.FollowersDatabase URI.")
	fs.StringVar(&following_database_uri, "following-database-uri", "", "A registered sfomuseum/go-activitypub/database.FollowingDatabase URI.")
	fs.StringVar(&notes_database_uri, "notes-database-uri", "", "A registered sfomuseum/go-activitypub/database.NotesDatabase URI.")
	fs.StringVar(&messages_database_uri, "messages-database-uri", "", "A registered sfomuseum/go-activitypub/database.MessagesDatabase URI.")
	fs.StringVar(&blocks_database_uri, "blocks-database-uri", "", "A registered sfomuseum/go-activitypub/database.BlocksDatabase URI.")
	fs.StringVar(&posts_database_uri, "posts-database-uri", "", "A registered sfomuseum/go-activitypub/database.PostsDatabase URI.")
	fs.StringVar(&post_tags_database_uri, "post-tags-database-uri", "", "A registered sfomuseum/go-activitypub/database.PostTagsDatabase URI.")
	fs.StringVar(&properties_database_uri, "properties-database-uri", "", "A registered sfomuseum/go-activitypub/database.PropertiesDatabase URI.")
	fs.StringVar(&likes_database_uri, "likes-database-uri", "", "A registered sfomuseum/go-activitypub/database.LikesDatabase URI.")
	fs.StringVar(&boosts_database_uri, "boosts-database-uri", "", "A registered sfomuseum/go-activitypub/database.BoostsDatabase URI.")

	fs.BoolVar(&allow_follow, "allow-follow", true, "Enable support for ActivityPub \"Follow\" activities.")
	fs.BoolVar(&allow_create, "allow-create", false, "Enable support for ActivityPub \"Create\" activities.")
	fs.BoolVar(&allow_likes, "allow-likes", true, "Enable support for ActivityPub \"Like\" activities.")
	fs.BoolVar(&allow_boosts, "allow-boosts", true, "Enable support for ActivityPub \"Announce\" (boost) activities.")
	fs.BoolVar(&allow_mentions, "allow-mentions", true, "If enabled allows posts (\"Create\" activities) to accounts not followed by author but where account is mentioned in post.")

	fs.StringVar(&server_uri, "server-uri", "http://localhost:8080", "A registered aaronland/go-http-server/server.Server URI.")
	fs.StringVar(&hostname, "hostname", "", "The hostname (domain) of the ActivityPub server delivering activities.")
	fs.BoolVar(&insecure, "insecure", false, "A boolean flag indicating the ActivityPub server delivering activities is insecure (not using TLS).")

	fs.StringVar(&process_message_queue_uri, "process-message-queue-uri", "null://", "A registered go-activitypub/queue.ProcessMessageQueue URI.")
	fs.StringVar(&process_follow_queue_uri, "process-follow-queue-uri", "null://", "A registered go-activitypub/queue.ProcessFollowQueue URI.")

	fs.BoolVar(&allow_remote_icon_uri, "allow-remote-icon-uri", false, "Allow account icons hosted on a remote host.")
	fs.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging.")
	fs.BoolVar(&disabled, "disabled", false, "Return a 503 Service unavailable response for all requests.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Start a HTTP (web) server to handle ActivityPub-related requests.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Valid options are:\n")
		fs.PrintDefaults()
	}

	return fs
}
