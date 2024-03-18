package followers

import (
	"context"
	"flag"
	"fmt"
	"log/slog"

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

	followers_db, err := activitypub.NewFollowersDatabase(ctx, opts.FollowersDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to instantiate followers database, %w", err)
	}

	defer followers_db.Close(ctx)

	followers := make(map[string]int)

	cb := func(ctx context.Context, addr string) error {

		count, exists := followers[addr]

		if !exists {
			count = 0
		}

		count += 1
		followers[addr] = count
		return nil
	}

	err = followers_db.GetAllFollowers(ctx, cb)

	count_followers := 0

	for addr, count := range followers {
		count_followers += 1
		fmt.Println(addr, count)
	}

	logger.Info("Followers", "count", count_followers)
	return nil
}
