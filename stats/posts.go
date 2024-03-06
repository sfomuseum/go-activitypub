package stats

import (
	"context"
	"fmt"
	"time"

	"github.com/sfomuseum/go-activitypub"
)

func CountPostsForDay(ctx context.Context, posts_db activitypub.PostsDatabase, day string) (int64, error) {

	t, err := time.Parse("2006-01-02", day)

	if err != nil {
		return 0, fmt.Errorf("Failed to parse day, %w", err)
	}

	start := t.Unix()
	end := start + ONEDAY

	return CountPostsForDateRange(ctx, posts_db, start, end)
}

func CountPostsForDateRange(ctx context.Context, posts_db activitypub.PostsDatabase, start int64, end int64) (int64, error) {

	count := int64(0)

	cb := func(ctx context.Context, id int64) error {
		count += 1
		return nil
	}

	err := posts_db.GetPostIdsForDateRange(ctx, start, end, cb)

	if err != nil {
		return 0, fmt.Errorf("Failed to count posts, %w", err)
	}

	return count, nil
}
