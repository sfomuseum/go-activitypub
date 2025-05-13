package create

import (
	"flag"
	"fmt"
	"os"

	"github.com/sfomuseum/go-flags/flagset"
)

var accounts_database_uri string
var activities_database_uri string
var followers_database_uri string
var posts_database_uri string
var post_tags_database_uri string
var deliveries_database_uri string

var delivery_queue_uri string

var account_name string
var message string
var in_reply_to string

var max_attempts int

var hostname string
var insecure bool
var verbose bool

var mode string
var lambda_function_uri string

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("follow")

	fs.StringVar(&mode, "mode", "cli", "Valid options are: cli, lambda and invoke.")
	fs.StringVar(&lambda_function_uri, "lambda-function-uri", "", "...")

	fs.StringVar(&accounts_database_uri, "accounts-database-uri", "null://", "A registered sfomuseum/go-activitypub/database.AccountsDatabase URI.")
	fs.StringVar(&activities_database_uri, "activities-database-uri", "null://", "A registered sfomuseum/go-activitypub/database.ActivitiesDatabase URI.")
	fs.StringVar(&followers_database_uri, "followers-database-uri", "null://", "A registered sfomuseum/go-activitypub/database.FollowersDatabase URI.")
	fs.StringVar(&posts_database_uri, "posts-database-uri", "null://", "A registered sfomuseum/go-activitypub/database.PostsDatabase URI.")
	fs.StringVar(&post_tags_database_uri, "post-tags-database-uri", "null://", "A registered sfomuseum/go-activitypub/database.PostTagsDatabase URI.")
	fs.StringVar(&deliveries_database_uri, "deliveries-database-uri", "null://", "A registered sfomuseum/go-activitypub/database.DeliveriesDatabase URI.")

	fs.StringVar(&delivery_queue_uri, "delivery-queue-uri", "synchronous://", "A registered sfomuseum/go-activitypub/queue/DeliveryQueue URI.")

	fs.StringVar(&account_name, "account-name", "", "The name of the go-activitypub account creating the post.")
	fs.StringVar(&hostname, "hostname", "localhost:8080", "The hostname (domain) of the ActivityPub server delivering activities.")
	fs.BoolVar(&insecure, "insecure", false, "A boolean flag indicating the ActivityPub server delivering activities is insecure (not using TLS).")

	fs.IntVar(&max_attempts, "max-attempts", 5, "The maximum number of attempts to deliver the activity.")
	fs.StringVar(&message, "message", "", "The body (content) of the message to post.")
	fs.StringVar(&in_reply_to, "in-reply-to", "", "The URI of that the post is in reply to (optional).")

	fs.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Create a new post (note, activity) on behalf of a registered go-activitypub account and schedule it for delivery to all their followers.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Valid options are:\n")
		fs.PrintDefaults()
	}

	return fs
}
