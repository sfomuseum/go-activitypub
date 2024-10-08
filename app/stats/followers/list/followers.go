package list

import (
	"context"
	"flag"
	"fmt"

	"github.com/sfomuseum/go-activitypub"
	"github.com/sfomuseum/go-activitypub/database"
	"github.com/sfomuseum/go-activitypub/slog"
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

	logger := slog.Logger()

	followers_db, err := database.NewFollowersDatabase(ctx, opts.FollowersDatabaseURI)

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
