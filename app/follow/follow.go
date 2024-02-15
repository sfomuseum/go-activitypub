package follow

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"

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

	db, err := activitypub.NewActorDatabase(ctx, opts.AccountDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create new database, %w", err)
	}

	acct, err := db.GetActor(ctx, opts.AccountId)

	if err != nil {
		return fmt.Errorf("Failed to retrieve account %s, %w", opts.AccountId, err)
	}

	for _, to_follow := range opts.Follow {

		acct_url, err := acct.ProfileURL(ctx, opts.URIs)

		if err != nil {
			return fmt.Errorf("Failed to derive profile URL for account, %w", err)
		}

		acct_url.Scheme = "http"

		follower := acct_url.String()

		follow_req, err := ap.NewFollowActivity(ctx, follower, to_follow)

		if err != nil {
			return fmt.Errorf("Failed to create follow activity, %w", err)
		}

		enc := json.NewEncoder(os.Stdout)
		enc.Encode(follow_req)
	}

	return nil
}
