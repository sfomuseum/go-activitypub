package stats

import (
	"context"
	"fmt"
	"time"

	"github.com/sfomuseum/go-activitypub"
)

func CountFollowingForDay(ctx context.Context, following_db activitypub.FollowingDatabase, day string) (int64, error) {

	t, err := time.Parse("2006-01-02", day)

	if err != nil {
		return 0, fmt.Errorf("Failed to parse day, %w", err)
	}

	start := t.Unix()
	end := start + ONEDAY

	return CountFollowingForDateRange(ctx, following_db, start, end)
}

func CountFollowingForDateRange(ctx context.Context, following_db activitypub.FollowingDatabase, start int64, end int64) (int64, error) {

	count := int64(0)

	cb := func(ctx context.Context, id int64) error {
		count += 1
		return nil
	}

	err := following_db.GetFollowingIdsForDateRange(ctx, start, end, cb)

	if err != nil {
		return 0, fmt.Errorf("Failed to count following, %w", err)
	}

	return count, nil
}
