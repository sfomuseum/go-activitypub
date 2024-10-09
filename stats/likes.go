package stats

import (
	"context"
	"fmt"

	"github.com/sfomuseum/go-activitypub/database"
)

func CountLikesForDateRange(ctx context.Context, likes_db database.LikesDatabase, start int64, end int64) (int64, error) {

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
