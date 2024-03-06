package stats

import (
	"context"
	"fmt"
	"time"

	"github.com/sfomuseum/go-activitypub"
)

func CountLikesForDay(ctx context.Context, likes_db activitypub.LikesDatabase, day string) (int64, error) {

	t, err := time.Parse("2006-01-02", day)

	if err != nil {
		return 0, fmt.Errorf("Failed to parse day, %w", err)
	}

	start := t.Unix()
	end := start + ONEDAY

	return CountLikesForDateRange(ctx, likes_db, start, end)
}

func CountLikesForDateRange(ctx context.Context, likes_db activitypub.LikesDatabase, start int64, end int64) (int64, error) {

	count := int64(0)

	cb := func(ctx context.Context, id int64) error {
		count += 1
		return nil
	}

	err := likes_db.GetLikeIdsForDateRange(ctx, start, end, cb)

	if err != nil {
		return 0, fmt.Errorf("Failed to count likes, %w", err)
	}

	return count, nil
}
