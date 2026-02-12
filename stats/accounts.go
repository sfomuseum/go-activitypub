package stats

import (
	"context"
	"fmt"

	"github.com/sfomuseum/go-activitypub/database"
)

func CountAccountsForDateRange(ctx context.Context, accounts_db database.AccountsDatabase, start int64, end int64) (int64, error) {

	count := int64(0)

	for _, err := range accounts_db.GetAccountIdsForDateRange(ctx, start, end) {

		if err != nil {
			return 0, fmt.Errorf("Failed to count accounts, %w", err)
		}

		count += 1
	}

	return count, nil
}
