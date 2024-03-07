package stats

import (
	"context"
	"fmt"

	"github.com/sfomuseum/go-activitypub"
)

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
