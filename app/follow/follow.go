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

	db, err := activitypub.NewAccountsDatabase(ctx, opts.AccountsDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create new database, %w", err)
	}

	// The person doing the following

	follower_acct, err := db.GetAccount(ctx, opts.AccountId)

	if err != nil {
		return fmt.Errorf("Failed to retrieve account %s, %w", opts.AccountId, err)
	}

	follower_id := follower_acct.Id
	following_id := opts.Follow

	follow_req, err := ap.NewFollowActivity(ctx, follower_id, following_id)

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
		logger.Info("Unfollowing successful")
	} else {
		logger.Info("Following successful")
	}

	return nil
}
