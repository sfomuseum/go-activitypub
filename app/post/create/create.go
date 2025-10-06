package create

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"

	aa_lambda "github.com/aaronland/go-aws/v3/lambda"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sfomuseum/go-activitypub"
	"github.com/sfomuseum/go-activitypub/database"
	"github.com/sfomuseum/go-activitypub/posts"
	"github.com/sfomuseum/go-activitypub/queue"
)

type Post struct {
	// The name of the go-activitypub account creating the post.
	AccountName string `json:"account_name"`
	// The body (content) of the message to post.
	Message string `json:"message"`
	// The URI of that the post is in reply to (optional).
	InReplyTo string `json:"in_reply_to,omitempty"`
}

func Run(ctx context.Context) error {
	fs := DefaultFlagSet()
	return RunWithFlagSet(ctx, fs)
}

func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet) error {

	opts, err := OptionsFromFlagSet(ctx, fs)

	if err != nil {
		return fmt.Errorf("Failed to derive options from flagset, %w", err)
	}

	return RunWithOptions(ctx, opts)
}

func RunWithOptions(ctx context.Context, opts *RunOptions) error {

	_, err := RunWithOptionsAndResponse(ctx, opts)
	return err
}

func RunWithOptionsAndResponse(ctx context.Context, opts *RunOptions) (string, error) {

	if opts.Verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Verbose logging enabled")
	}

	accounts_db, err := database.NewAccountsDatabase(ctx, opts.AccountsDatabaseURI)

	if err != nil {
		return "", fmt.Errorf("Failed to create accounts database, %w", err)
	}

	defer accounts_db.Close(ctx)

	activities_db, err := database.NewActivitiesDatabase(ctx, opts.ActivitiesDatabaseURI)

	if err != nil {
		return "", fmt.Errorf("Failed to create activities database, %w", err)
	}

	defer activities_db.Close(ctx)

	followers_db, err := database.NewFollowersDatabase(ctx, opts.FollowersDatabaseURI)

	if err != nil {
		return "", fmt.Errorf("Failed to instantiate followers database, %w", err)
	}

	defer followers_db.Close(ctx)

	posts_db, err := database.NewPostsDatabase(ctx, opts.PostsDatabaseURI)

	if err != nil {
		return "", fmt.Errorf("Failed to create instantiate posts database, %w", err)
	}

	defer posts_db.Close(ctx)

	post_tags_db, err := database.NewPostTagsDatabase(ctx, opts.PostTagsDatabaseURI)

	if err != nil {
		return "", fmt.Errorf("Failed to create instantiate post tags database, %w", err)
	}

	defer post_tags_db.Close(ctx)

	deliveries_db, err := database.NewDeliveriesDatabase(ctx, opts.DeliveriesDatabaseURI)

	if err != nil {
		return "", fmt.Errorf("Failed to create instantiate deliveries database, %w", err)
	}

	defer deliveries_db.Close(ctx)

	delivery_q, err := queue.NewDeliveryQueue(ctx, opts.DeliveryQueueURI)

	if err != nil {
		return "", fmt.Errorf("Failed to create new delivery queue, %w", err)
	}

	run := func(ctx context.Context, opts *RunOptions) (string, error) {

		message := opts.Message

		if message == "" {
			return "", fmt.Errorf("Empty message string")
		}

		logger := slog.Default()
		logger = logger.With("account", opts.AccountName)

		acct, err := accounts_db.GetAccountWithName(ctx, opts.AccountName)

		if err != nil {
			return "", fmt.Errorf("Failed to retrieve account %s, %w", opts.AccountName, err)
		}

		logger = logger.With("account id", acct.Id)

		post_opts := &posts.AddPostOptions{
			URIs:          opts.URIs,
			PostsDatabase: posts_db,
			// aka mentions
			PostTagsDatabase: post_tags_db,
		}

		logger.Debug("Add post", "message", message)

		post, mentions, err := posts.AddPost(ctx, post_opts, acct, opts.Message)

		if err != nil {
			return "", fmt.Errorf("Failed to add post, %w", err)
		}

		if opts.InReplyTo != "" {
			post.InReplyTo = opts.InReplyTo
		}

		logger = logger.With("post id", post.Id)

		ap_activity, err := posts.ActivityFromPost(ctx, opts.URIs, acct, post, mentions)

		if err != nil {
			return "", fmt.Errorf("Failed to create new (create) activity, %w", err)
		}

		activity, err := activitypub.NewActivity(ctx, ap_activity)

		if err != nil {
			return "", fmt.Errorf("Failed to create new AP wrapper, %w", err)
		}

		activity.ActivityType = activitypub.PostActivityType
		activity.ActivityTypeId = post.Id
		activity.AccountId = acct.Id

		err = activities_db.AddActivity(ctx, activity)

		if err != nil {
			return "", fmt.Errorf("Failed to add activity, %w", err)
		}

		// TBD add activity to activities_db

		logger = logger.With("activity id", activity.Id)

		deliver_opts := &queue.DeliverActivityToFollowersOptions{
			AccountsDatabase:   accounts_db,
			FollowersDatabase:  followers_db,
			DeliveriesDatabase: deliveries_db,
			DeliveryQueue:      delivery_q,
			Activity:           activity,
			Mentions:           mentions,
			URIs:               opts.URIs,
			MaxAttempts:        opts.MaxAttempts,
		}

		logger.Debug("Deliver activity")

		err = queue.DeliverActivityToFollowers(ctx, deliver_opts)

		if err != nil {
			return "", fmt.Errorf("Failed to deliver post, %w", err)
		}

		post_url := acct.PostURL(ctx, opts.URIs, post).String()
		return post_url, nil
	}

	switch opts.Mode {
	case "cli":

		message := opts.Message

		if message == "-" {

			scanner := bufio.NewScanner(os.Stdin)

			for scanner.Scan() {
				message = fmt.Sprintf("%s %s", message, scanner.Text())
			}

			if scanner.Err() != nil {
				return "", fmt.Errorf("Failed to scan input, %w", err)
			}

			opts.Message = message
		}

		post_url, err := run(ctx, opts)

		if err != nil {
			return "", fmt.Errorf("Failed to post message, %w", err)
		}

		slog.Info("Delivered post", "post url", post_url)
		return post_url, err

	case "lambda":

		handle := func(ctx context.Context, post *Post) (string, error) {

			opts.AccountName = post.AccountName
			opts.Message = post.Message
			opts.InReplyTo = post.InReplyTo

			return run(ctx, opts)
		}

		lambda.Start(handle)
		return "", nil

	case "invoke":

		post := &Post{
			AccountName: opts.AccountName,
			Message:     opts.Message,
		}

		if opts.InReplyTo != "" {
			post.InReplyTo = opts.InReplyTo
		}

		fn, err := aa_lambda.NewLambdaFunction(ctx, opts.LambdaFunctionURI)

		if err != nil {
			return "", fmt.Errorf("Failed to create new lambda function, %w", err)
		}

		rsp, err := fn.Invoke(ctx, post)

		if err != nil {
			slog.Error("Failed to invoke lambda function", "error", err)
			return "", fmt.Errorf("Failed to invoke Lambda function, %w", err)
		}

		post_url := string(rsp.Payload)
		post_url = strings.TrimLeft(post_url, `"`)
		post_url = strings.TrimRight(post_url, `"`)

		slog.Info("Delivered post", "post url", post_url)
		return post_url, nil

	default:
		return "", fmt.Errorf("Invalid or unsupported mode")
	}

	return "", nil
}
