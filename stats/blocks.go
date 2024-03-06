package stats

import (
	"context"
	"fmt"
	"time"

	"github.com/sfomuseum/go-activitypub"
)

func CountBlocksForDay(ctx context.Context, blocks_db activitypub.BlocksDatabase, day string) (int64, error) {

	t, err := time.Parse("2006-01-02", day)

	if err != nil {
		return 0, fmt.Errorf("Failed to parse day, %w", err)
	}

	start := t.Unix()
	end := start + ONEDAY

	return CountBlocksForDateRange(ctx, blocks_db, start, end)
}

func CountBlocksForDateRange(ctx context.Context, blocks_db activitypub.BlocksDatabase, start int64, end int64) (int64, error) {

	count := int64(0)

	cb := func(ctx context.Context, id int64) error {
		count += 1
		return nil
	}

	err := blocks_db.GetBlockIdsForDateRange(ctx, start, end, cb)

	if err != nil {
		return 0, fmt.Errorf("Failed to count blocks, %w", err)
	}

	return count, nil
}
