package stats

import (
	"context"
	"fmt"

	"github.com/sfomuseum/go-activitypub"
)

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
