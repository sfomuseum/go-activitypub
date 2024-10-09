package blocks

import (
	"context"
	"fmt"

	"github.com/sfomuseum/go-activitypub"
	"github.com/sfomuseum/go-activitypub/database"
)

func IsBlockedByAccount(ctx context.Context, db database.BlocksDatabase, account_id int64, host string, name string) (bool, error) {

	_, err := db.GetBlockWithAccountIdAndAddress(ctx, account_id, host, name)

	if err == nil {
		return true, nil
	}

	if err != activitypub.ErrNotFound {
		return false, fmt.Errorf("Failed to retrieve block with account and address, %w", err)
	}

	if name == "*" {
		return false, nil
	}

	return IsBlockedByAccount(ctx, db, account_id, host, "*")
}
