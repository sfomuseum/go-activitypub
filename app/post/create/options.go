package create

import (
	"context"
	"flag"
	"fmt"

	"github.com/sfomuseum/go-activitypub/uris"
	"github.com/sfomuseum/go-flags/flagset"
)

type RunOptions struct {
	AccountsDatabaseURI   string
	ActivitiesDatabaseURI string
	FollowersDatabaseURI  string
	PostsDatabaseURI      string
	PostTagsDatabaseURI   string
	DeliveriesDatabaseURI string
	DeliveryQueueURI      string
	AccountName           string
	Message               string
	InReplyTo             string
	MaxAttempts           int
	URIs                  *uris.URIs
	Mode                  string
	LambdaFunctionURI     string
	Verbose               bool
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
