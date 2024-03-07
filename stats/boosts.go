package stats

import (
	"context"
	"fmt"

	"github.com/sfomuseum/go-activitypub"
)

func CountBoostsForDateRange(ctx context.Context, boosts_db activitypub.BoostsDatabase, start int64, end int64) (int64, error) {

	count := int64(0)

	cb := func(ctx context.Context, id int64) error {
		count += 1
		return nil
	}

	err := boosts_db.GetBoostIdsForDateRange(ctx, start, end, cb)

	if err != nil {
		return 0, fmt.Errorf("Failed to count boosts, %w", err)
	}

	return count, nil
}
