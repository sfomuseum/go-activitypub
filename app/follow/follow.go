package follow

import (
	"context"
	"flag"
	"fmt"
	"log/slog"

	"github.com/sfomuseum/go-activitypub"
	"github.com/sfomuseum/go-activitypub/ap"
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
		return fmt.Errorf("Failed to initialize accounts database, %w", err)
	}

	following_db, err := activitypub.NewFollowingDatabase(ctx, opts.FollowingDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to initialize following database, %w", err)
	}

	follower_acct, err := accounts_db.GetAccount(ctx, opts.AccountId)

	if err != nil {
		return fmt.Errorf("Failed to retrieve account %s, %w", opts.AccountId, err)
	}

	follower_id := follower_acct.Id
	following_id := opts.Follow

	// See this? It is important to pass the fully-qualifier follower URI so the
	// endpoint receiving the follow activity can figure out where (which hostname)
	// to make a webfinger/profile query.

	follower_uri := fmt.Sprintf("%s@%s", follower_id, opts.Hostname)

	follow_req, err := ap.NewFollowActivity(ctx, follower_uri, following_id)

	if err != nil {
		return fmt.Errorf("Failed to create follow activity, %w", err)
	}

	if opts.Undo {
		follow_req.Type = "Undo"
	}

	post_opts := &activitypub.PostToAccountOptions{
		From:     follower_acct,
		To:       following_id,
		Hostname: opts.Hostname,
		URIs:     opts.URIs,
		Message:  follow_req,
	}

	_, err = activitypub.PostToAccount(ctx, post_opts)

	if undo {

		err := following_db.UnFollow(ctx, follower_id, following_id)

		if err != nil {
			return fmt.Errorf("Unfollow request was successful but unable to register unfollowing locally, %w", err)
		}

		logger.Info("Unfollowing successful")
		return nil
	}

	err = following_db.Follow(ctx, follower_id, following_id)

	if err != nil {
		return fmt.Errorf("Follow request was successful but unable to register following locally, %w", err)
	}

	logger.Info("Following successful")
	return nil
}
