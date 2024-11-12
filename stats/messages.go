package stats

import (
	"context"
	"fmt"

	"github.com/sfomuseum/go-activitypub/database"
)

func CountMessagesForDateRange(ctx context.Context, messages_db database.MessagesDatabase, start int64, end int64) (int64, error) {

	count := int64(0)

	cb := func(ctx context.Context, id int64) error {
		count += 1
		return nil
	}

	err := messages_db.GetMessageIdsForDateRange(ctx, start, end, cb)

	if err != nil {
		return 0, fmt.Errorf("Failed to count messages, %w", err)
	}

	return count, nil
}
