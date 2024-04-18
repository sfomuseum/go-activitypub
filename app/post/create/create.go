package create

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/sfomuseum/go-activitypub"
	ap_slog "github.com/sfomuseum/go-activitypub/slog"
)

func Run(ctx context.Context, logger *slog.Logger) error {
	fs := DefaultFlagSet()
	return RunWithFlagSet(ctx, fs, logger)
}

func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet, logger *slog.Logger) error {

	opts, err := OptionsFromFlagSet(ctx, fs)

	if err != nil {
		return fmt.Errorf("Failed to derive options from flagset, %w", err)
	}

	return RunWithOptions(ctx, opts, logger)
}

func RunWithOptions(ctx context.Context, opts *RunOptions, logger *slog.Logger) error {

	ap_slog.ConfigureLogger(logger, opts.Verbose)

	accounts_db, err := activitypub.NewAccountsDatabase(ctx, opts.AccountsDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create new database, %w", err)
	}

	defer accounts_db.Close(ctx)

	followers_db, err := activitypub.NewFollowersDatabase(ctx, opts.FollowersDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to instantiate followers database, %w", err)
	}

	defer followers_db.Close(ctx)

	posts_db, err := activitypub.NewPostsDatabase(ctx, opts.PostsDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create instantiate posts database, %w", err)
	}

	defer posts_db.Close(ctx)

	post_tags_db, err := activitypub.NewPostTagsDatabase(ctx, opts.PostTagsDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create instantiate post tags database, %w", err)
	}

	defer post_tags_db.Close(ctx)

	deliveries_db, err := activitypub.NewDeliveriesDatabase(ctx, opts.DeliveriesDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create instantiate deliveries database, %w", err)
	}

	defer deliveries_db.Close(ctx)

	delivery_q, err := activitypub.NewDeliveryQueue(ctx, opts.DeliveryQueueURI)

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

	acct, err := accounts_db.GetAccountWithName(ctx, opts.AccountName)

	if err != nil {
		return fmt.Errorf("Failed to retrieve account %s, %w", opts.AccountName, err)
	}

	post_opts := &activitypub.AddPostOptions{
		URIs:             opts.URIs,
		PostsDatabase:    posts_db,
		PostTagsDatabase: post_tags_db,
	}

	post, post_tags, err := activitypub.AddPost(ctx, post_opts, acct, opts.Message)

	if err != nil {
		return fmt.Errorf("Failed to add post, %w", err)
	}

	if opts.InReplyTo != "" {
		post.InReplyTo = opts.InReplyTo
	}

	deliver_opts := &activitypub.DeliverPostToFollowersOptions{
		AccountsDatabase:   accounts_db,
		FollowersDatabase:  followers_db,
		PostTagsDatabase:   post_tags_db,
		DeliveriesDatabase: deliveries_db,
		DeliveryQueue:      delivery_q,
		Post:               post,
		PostTags:           post_tags,
		URIs:               opts.URIs,
		MaxAttempts:        opts.MaxAttempts,
	}

	err = activitypub.DeliverPostToFollowers(ctx, deliver_opts)

	if err != nil {
		return fmt.Errorf("Failed to deliver post, %w", err)
	}

	logger.Info("Delivered post", "ID", acct.PostURL(ctx, opts.URIs, post).String())
	return nil
}
