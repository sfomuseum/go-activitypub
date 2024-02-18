package post

import (
	"context"
	"flag"
	"fmt"
	"log/slog"

	"github.com/sfomuseum/go-activitypub"
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

	slog.SetDefault(logger)

	accounts_db, err := activitypub.NewAccountsDatabase(ctx, opts.AccountsDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create new database, %w", err)
	}

	followers_db, err := activitypub.NewFollowersDatabase(ctx, opts.FollowersDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to instantiate followers database, %w", err)
	}

	posts_db, err := activitypub.NewPostsDatabase(ctx, opts.PostsDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create instatiate posts database, %w", err)
	}

	delivery_q, err := activitypub.NewDeliveryQueue(ctx, opts.DeliveryQueueURI)

	if err != nil {
		return fmt.Errorf("Failed to create new delivery queue, %w", err)
	}

	acct, err := accounts_db.GetAccount(ctx, opts.AccountId)

	if err != nil {
		return fmt.Errorf("Failed to retrieve account %s, %w", opts.AccountId, err)
	}

	p, err := activitypub.NewPost(ctx, acct, []byte("Hello world"))

	if err != nil {
		return fmt.Errorf("Failed to create new post, %w", err)
	}

	err = posts_db.AddPost(ctx, p)

	if err != nil {
		return fmt.Errorf("Failed to add post, %w", err)
	}

	deliver_opts := &activitypub.DeliverPostToFollowersOptions{
		AccountsDatabase:  accounts_db,
		FollowersDatabase: followers_db,
		DeliveryQueue:     delivery_q,
		Post:              p,
		Hostname:          opts.Hostname,
		URIs:              opts.URIs,
	}

	err = activitypub.DeliverPostToFollowers(ctx, deliver_opts)

	if err != nil {
		return fmt.Errorf("Failed to deliver post, %w", err)
	}

	return nil
}
