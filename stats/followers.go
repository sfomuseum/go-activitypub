package stats

import (
	"context"
	"fmt"
	"time"

	"github.com/sfomuseum/go-activitypub"
)

func CountFollowersForDay(ctx context.Context, followers_db activitypub.FollowersDatabase, day string) (int64, error) {

	t, err := time.Parse("2006-01-02", day)

	if err != nil {
		return 0, fmt.Errorf("Failed to parse day, %w", err)
	}

	start := t.Unix()
	end := start + ONEDAY

	return CountFollowersForDateRange(ctx, followers_db, start, end)
}

func CountFollowersForDateRange(ctx context.Context, followers_db activitypub.FollowersDatabase, start int64, end int64) (int64, error) {

	count := int64(0)

	cb := func(ctx context.Context, id int64) error {
		count += 1
		return nil
	}

	err := followers_db.GetFollowerIdsForDateRange(ctx, start, end, cb)

	if err != nil {
		return 0, fmt.Errorf("Failed to count followers, %w", err)
	}

	return count, nil
}
