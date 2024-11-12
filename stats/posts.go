package stats

import (
	"context"
	"fmt"

	"github.com/sfomuseum/go-activitypub/database"
)

func CountPostsForDateRange(ctx context.Context, posts_db database.PostsDatabase, start int64, end int64) (int64, error) {

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
