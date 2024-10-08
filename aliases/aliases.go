package aliases

import (
	"context"
	"fmt"

	"github.com/sfomuseum/go-activitypub/database"
)

func IsAliasNameTaken(ctx context.Context, aliases_db database.AliasesDatabase, name string) (bool, error) {

	_, err := aliases_db.GetAliasWithName(ctx, name)

	if err != nil {

		if err != ErrNotFound {
			return false, fmt.Errorf("Failed to determine is name is taken, %w", err)
		}

		return false, nil
	}

	return true, nil
}
