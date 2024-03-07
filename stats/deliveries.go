package stats

import (
	"context"
	"fmt"

	"github.com/sfomuseum/go-activitypub"
)

func CountDeliveriesForDateRange(ctx context.Context, deliveries_db activitypub.DeliveriesDatabase, start int64, end int64) (int64, error) {

	count := int64(0)

	cb := func(ctx context.Context, id int64) error {
		count += 1
		return nil
	}

	err := deliveries_db.GetDeliveryIdsForDateRange(ctx, start, end, cb)

	if err != nil {
		return 0, fmt.Errorf("Failed to count deliveries, %w", err)
	}

	return count, nil
}
