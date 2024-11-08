package deliver

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

var subscriber_uri string

var delivery_queue_uri string
var max_attempts int

var activity_id int64
var mode string

// Allows posts to accounts not followed by author but where account is mentioned in post
var allow_mentions bool

var hostname string
var insecure bool
var verbose bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("follow")

	fs.StringVar(&accounts_database_uri, "accounts-database-uri", "", "A registered sfomuseum/go-activitypub/database.AccountsDatabase URI.")
	fs.StringVar(&activities_database_uri, "activities-database-uri", "", "A registered sfomuseum/go-activitypub/database.ActivitiesDatabase URI.")
	fs.StringVar(&followers_database_uri, "followers-database-uri", "", "A registered sfomuseum/go-activitypub/database.FollowersDatabase URI.")
	fs.StringVar(&posts_database_uri, "posts-database-uri", "", "A registered sfomuseum/go-activitypub/database.PostsDatabase URI.")
	fs.StringVar(&post_tags_database_uri, "post-tags-database-uri", "null://", "A registered sfomuseum/go-activitypub/database.PostTagsDatabase URI.")
	fs.StringVar(&deliveries_database_uri, "deliveries-database-uri", "", "A registered sfomuseum/go-activitypub/database.DeliveriesDatabase URI.")

	fs.BoolVar(&allow_mentions, "allow-mentions", true, "Enable support for processing mentions in (post) activities. This enabled posts to accounts not followed by author but where account is mentioned in post.")
	fs.StringVar(&delivery_queue_uri, "delivery-queue-uri", "synchronous://", "A registered sfomuseum/go-activitypub/queue.DeliveryQueue URI.")

	fs.IntVar(&max_attempts, "max-attempts", 5, "The maximum number of attempts to deliver the activity.")
	// fs.Int64Var(&activity_id, "activity-id", "The unique sfomuseum/go-activitypub.Activity ID to deliver.")

	fs.StringVar(&subscriber_uri, "subscriber-uri", "", "A valid sfomuseum/go-pubsub/subscriber URI. Required if -mode parameter is 'pubsub'.")

	fs.StringVar(&mode, "mode", "", "The operation mode for delivering activities. Valid options are: lambda, pubsub. \"cli\" mode is currently disabled.")

	fs.StringVar(&hostname, "hostname", "localhost:8080", "The hostname of the ActivityPub server delivering activities.")
	fs.BoolVar(&insecure, "insecure", false, "A boolean flag indicating the ActivityPub server delivering activities is insecure (not using TLS).")
	fs.BoolVar(&verbose, "verbose", false, "Enable verbose logging")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Deliver an ActivityPub activity to subscribers.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Valid options are:\n")
		fs.PrintDefaults()
	}

	return fs
}
