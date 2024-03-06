package stats

import (
	"context"
	"fmt"
	"time"

	"github.com/sfomuseum/go-activitypub"
)

func CountAccountsForDay(ctx context.Context, accounts_db activitypub.AccountsDatabase, day string) (int64, error) {

	t, err := time.Parse("2006-01-02", day)

	if err != nil {
		return 0, fmt.Errorf("Failed to parse day, %w", err)
	}

	start := t.Unix()
	end := start + ONEDAY

	return CountAccountsForDateRange(ctx, accounts_db, start, end)
}

func CountAccountsForDateRange(ctx context.Context, accounts_db activitypub.AccountsDatabase, start int64, end int64) (int64, error) {

	count := int64(0)

	cb := func(ctx context.Context, id int64) error {
		count += 1
		return nil
	}

	err := accounts_db.GetAccountIdsForDateRange(ctx, start, end, cb)

	if err != nil {
		return 0, fmt.Errorf("Failed to count accounts, %w", err)
	}

	return count, nil
}
