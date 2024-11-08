package database

import (
	"context"

	"github.com/sfomuseum/go-activitypub"
)

type NullAccountsDatabase struct {
	AccountsDatabase
}

func init() {
	ctx := context.Background()
	err := RegisterAccountsDatabase(ctx, "null", NewNullAccountsDatabase)

	if err != nil {
		panic(err)
	}
}

func NewNullAccountsDatabase(ctx context.Context, uri string) (AccountsDatabase, error) {
	db := &NullAccountsDatabase{}
	return db, nil
}

func (db *NullAccountsDatabase) GetAccounts(ctx context.Context, cb GetAccountsCallbackFunc) error {
	return nil
}

func (db *NullAccountsDatabase) GetAccountIdsForDateRange(ctx context.Context, start int64, end int64, cb GetAccountIdsCallbackFunc) error {
	return nil
}

func (db *NullAccountsDatabase) AddAccount(ctx context.Context, a *activitypub.Account) error {
	return nil
}

func (db *NullAccountsDatabase) GetAccountWithId(ctx context.Context, id int64) (*activitypub.Account, error) {
	return nil, activitypub.ErrNotFound
}

func (db *NullAccountsDatabase) GetAccountWithName(ctx context.Context, name string) (*activitypub.Account, error) {
	return nil, activitypub.ErrNotFound
}

func (db *NullAccountsDatabase) UpdateAccount(ctx context.Context, acct *activitypub.Account) error {
	return nil
}

func (db *NullAccountsDatabase) RemoveAccount(ctx context.Context, acct *activitypub.Account) error {
	return nil
}

func (db *NullAccountsDatabase) Close(ctx context.Context) error {
	return nil
}
