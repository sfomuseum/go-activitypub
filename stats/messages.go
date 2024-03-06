package stats

import (
	"context"
	"fmt"
	"time"

	"github.com/sfomuseum/go-activitypub"
)

func CountMessagesForDay(ctx context.Context, messages_db activitypub.MessagesDatabase, day string) (int64, error) {

	t, err := time.Parse("2006-01-02", day)

	if err != nil {
		return 0, fmt.Errorf("Failed to parse day, %w", err)
	}

	start := t.Unix()
	end := start + ONEDAY

	return CountMessagesForDateRange(ctx, messages_db, start, end)
}

func CountMessagesForDateRange(ctx context.Context, messages_db activitypub.MessagesDatabase, start int64, end int64) (int64, error) {

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
