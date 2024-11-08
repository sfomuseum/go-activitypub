package example

import (
	"flag"
	"fmt"
	"os"

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

	fs.StringVar(&mode, "mode", "cli", "The mode of operation. Valid options are: cli, lambda.")
	fs.Var(&message_ids, "id", "One or more message IDs to process (if -mode=cli).")

	fs.StringVar(&accounts_database_uri, "accounts-database-uri", "", "A registered sfomuseum/go-activitypub/database.AccountsDatabase URI.")
	fs.StringVar(&properties_database_uri, "properties-database-uri", "", "A registered sfomuseum/go-activitypub/database.PropertiesDatabase URI.")
	fs.StringVar(&activities_database_uri, "activities-database-uri", "", "A registered sfomuseum/go-activitypub/database.ActivitiesDatabase URI.")
	fs.StringVar(&messages_database_uri, "messages-database-uri", "", "A registered sfomuseum/go-activitypub/database.MessagesDatabase URI.")
	fs.StringVar(&notes_database_uri, "notes-database-uri", "", "A registered sfomuseum/go-activitypub/database.NotesDatabase URI.")
	fs.StringVar(&posts_database_uri, "posts-database-uri", "", "A registered sfomuseum/go-activitypub/database.PostsDatabase URI.")
	fs.StringVar(&post_tags_database_uri, "post-tags-database-uri", "", "A registered sfomuseum/go-activitypub/database.PostTagsDatabase URI.")
	fs.StringVar(&deliveries_database_uri, "deliveries-database-uri", "", "A registered sfomuseum/go-activitypub/database.DeliveriesDatabase URI.")
	fs.StringVar(&followers_database_uri, "followers-database-uri", "", "A registered sfomuseum/go-activitypub/database.FollowersDatabase URI.")

	fs.StringVar(&delivery_queue_uri, "delivery-queue-uri", "synchronous://", "A registered sfomuseum/go-activitypub/queue.DeliveryQueue URI.")

	fs.IntVar(&max_attempts, "max-attempts", 5, "The maximum number of attempts to try delivering activities.")
	fs.StringVar(&hostname, "hostname", "", "The hostname of the ActivityPub server delivering activities.")
	fs.BoolVar(&insecure, "insecure", false, "A boolean flag indicating the ActivityPub server delivering activities is insecure (not using TLS).")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "An example application for processing messages delivered through a go-activitypub/queue.ProcessMessageQueue publisher.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Valid options are:\n")
		fs.PrintDefaults()
	}

	return fs
}
