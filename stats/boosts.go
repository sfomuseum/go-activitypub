package stats

import (
	"context"
	"fmt"
	"time"

	"github.com/sfomuseum/go-activitypub"
)

func CountBoostsForDay(ctx context.Context, boosts_db activitypub.BoostsDatabase, day string) (int64, error) {

	t, err := time.Parse("2006-01-02", day)

	if err != nil {
		return 0, fmt.Errorf("Failed to parse day, %w", err)
	}

	start := t.Unix()
	end := start + ONEDAY

	return CountBoostsForDateRange(ctx, boosts_db, start, end)
}

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
