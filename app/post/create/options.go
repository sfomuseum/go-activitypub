package create

import (
	"context"
	"flag"
	"fmt"

	"github.com/sfomuseum/go-activitypub/uris"
	"github.com/sfomuseum/go-flags/flagset"
)

type RunOptions struct {
	// A registered sfomuseum/go-activitypub/database.AccountsDatabase URI.
	AccountsDatabaseURI string
	// A registered sfomuseum/go-activitypub/database.ActivitiesDatabase URI.
	ActivitiesDatabaseURI string
	// A registered sfomuseum/go-activitypub/database.FollowersDatabase URI.
	FollowersDatabaseURI string
	// A registered sfomuseum/go-activitypub/database.PostsDatabase URI.
	PostsDatabaseURI string
	// A registered sfomuseum/go-activitypub/database.PostTagsDatabase URI.
	PostTagsDatabaseURI string
	// A registered sfomuseum/go-activitypub/database.DeliveriesDatabase URI.
	DeliveriesDatabaseURI string
	// A registered sfomuseum/go-activitypub/queue/DeliveryQueue URI.
	DeliveryQueueURI string
	// The name of the go-activitypub account creating the post.
	AccountName string
	// The body (content) of the message to post.
	Message string
	// The URI of that the post is in reply to (optional).
	InReplyTo string
	// The maximum number of attempts to deliver the activity.
	MaxAttempts int
	// The operating mode for creating new posts. Valid options are: cli, lambda and invoke, where "lambda"
	// means to run as an AWS Lambda function and "invoke" means to invoke this tool as a specific Lambda function.
	Mode string
	// A valid aaronland/go-aws-lambda.LambdaFunction URI in the form of "lambda://FUNCTION_NAME}?region={AWS_REGION}&credentials={CREDENTIALS}".
	// This flag is required if the -mode flag is "invoke".
	LambdaFunctionURI string
	// Enable verbose (debug) logging.
	Verbose bool
	URIs    *uris.URIs
}

func OptionsFromFlagSet(ctx context.Context, fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	err := flagset.SetFlagsFromEnvVars(fs, "ACTIVITYPUB")

	if err != nil {
		return nil, fmt.Errorf("Failed to derive flags from environment variables, %w", err)
	}

	if mode == "invoke" && lambda_function_uri == "" {
		return nil, fmt.Errorf("Empty -lambda-function-uri flag")
	}

	uris_table := uris.DefaultURIs()
	uris_table.Hostname = hostname
	uris_table.Insecure = insecure

	opts := &RunOptions{
		AccountsDatabaseURI:   accounts_database_uri,
		ActivitiesDatabaseURI: activities_database_uri,
		FollowersDatabaseURI:  followers_database_uri,
		PostsDatabaseURI:      posts_database_uri,
		PostTagsDatabaseURI:   post_tags_database_uri,
		DeliveriesDatabaseURI: deliveries_database_uri,
		DeliveryQueueURI:      delivery_queue_uri,
		AccountName:           account_name,
		Message:               message,
		InReplyTo:             in_reply_to,
		URIs:                  uris_table,
		Verbose:               verbose,
		Mode:                  mode,
		LambdaFunctionURI:     lambda_function_uri,
		MaxAttempts:           max_attempts,
	}

	return opts, nil
}
