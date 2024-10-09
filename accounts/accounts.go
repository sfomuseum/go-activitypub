package accounts

import (
	"context"
	"fmt"
	"time"

	"github.com/sfomuseum/go-activitypub"
	"github.com/sfomuseum/go-activitypub/database"
)

func AddAccount(ctx context.Context, db database.AccountsDatabase, a *activitypub.Account) (*activitypub.Account, error) {

	now := time.Now()
	ts := now.Unix()

	a.Created = ts
	a.LastModified = ts

	err := db.AddAccount(ctx, a)

	if err != nil {
		return nil, fmt.Errorf("Failed to add account, %w", err)
	}

	return a, nil
}

func UpdateAccount(ctx context.Context, db database.AccountsDatabase, a *activitypub.Account) (*activitypub.Account, error) {

	now := time.Now()
	ts := now.Unix()

	a.LastModified = ts

	err := db.UpdateAccount(ctx, a)

	if err != nil {
		return nil, fmt.Errorf("Failed to update account, %w", err)
	}

	return a, nil
}

func IsAccountNameTaken(ctx context.Context, db database.AccountsDatabase, name string) (bool, error) {

	_, err := db.GetAccountWithName(ctx, name)

	if err != nil {

		if err != activitypub.ErrNotFound {
			return false, fmt.Errorf("Failed to determine is name is taken, %w", err)
		}

		return false, nil
	}

	return true, nil
}
