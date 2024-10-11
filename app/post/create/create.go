package create

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	// "github.com/sfomuseum/go-activitypub/ap"
	"github.com/sfomuseum/go-activitypub/database"
	"github.com/sfomuseum/go-activitypub/posts"
	"github.com/sfomuseum/go-activitypub/queue"
)

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

	if opts.Verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Verbose logging enabled")
	}

	logger := slog.Default()

	accounts_db, err := database.NewAccountsDatabase(ctx, opts.AccountsDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create new database, %w", err)
	}

	defer accounts_db.Close(ctx)

	followers_db, err := database.NewFollowersDatabase(ctx, opts.FollowersDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to instantiate followers database, %w", err)
	}

	defer followers_db.Close(ctx)

	posts_db, err := database.NewPostsDatabase(ctx, opts.PostsDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create instantiate posts database, %w", err)
	}

	defer posts_db.Close(ctx)

	post_tags_db, err := database.NewPostTagsDatabase(ctx, opts.PostTagsDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create instantiate post tags database, %w", err)
	}

	defer post_tags_db.Close(ctx)

	deliveries_db, err := database.NewDeliveriesDatabase(ctx, opts.DeliveriesDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create instantiate deliveries database, %w", err)
	}

	defer deliveries_db.Close(ctx)

	delivery_q, err := queue.NewDeliveryQueue(ctx, opts.DeliveryQueueURI)

	if err != nil {
		return fmt.Errorf("Failed to create new delivery queue, %w", err)
	}

	message := opts.Message

	if message == "-" {

		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() {
			message = fmt.Sprintf("%s %s", message, scanner.Text())
		}

		if scanner.Err() != nil {
			return fmt.Errorf("Failed to scan input, %w", err)
		}
	}

	if message == "" {
		return fmt.Errorf("Empty message string")
	}

	logger = logger.With("account", opts.AccountName)

	acct, err := accounts_db.GetAccountWithName(ctx, opts.AccountName)

	if err != nil {
		return fmt.Errorf("Failed to retrieve account %s, %w", opts.AccountName, err)
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
		return fmt.Errorf("Failed to add post, %w", err)
	}

	if opts.InReplyTo != "" {
		post.InReplyTo = opts.InReplyTo
	}

	logger = logger.With("post id", post.Id)

	activity, err := posts.ActivityFromPost(ctx, opts.URIs, acct, post, mentions)

	if err != nil {
		return fmt.Errorf("Failed to create new (create) activity, %w", err)
	}

	logger = logger.With("activity id", activity.Id)

	deliver_opts := &queue.DeliverActivityToFollowersOptions{
		AccountsDatabase:   accounts_db,
		FollowersDatabase:  followers_db,
		DeliveriesDatabase: deliveries_db,
		DeliveryQueue:      delivery_q,
		Activity:           activity,
		PostId:             post.Id,
		Mentions:           mentions,
		URIs:               opts.URIs,
		MaxAttempts:        opts.MaxAttempts,
	}

	logger.Debug("Deliver activity")

	err = queue.DeliverActivityToFollowers(ctx, deliver_opts)

	if err != nil {
		return fmt.Errorf("Failed to deliver post, %w", err)
	}

	logger.Info("Delivered post", "post url", acct.PostURL(ctx, opts.URIs, post).String())
	return nil
}
