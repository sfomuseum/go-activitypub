package activitypub

import (
	"context"
)

type NullAccountsDatabase struct {
	AccountsDatabase
}

func init() {
	ctx := context.Background()
	RegisterAccountsDatabase(ctx, "null", NewNullAccountsDatabase)

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

func (db *NullAccountsDatabase) AddAccount(ctx context.Context, a *Account) error {
	return nil
}

func (db *NullAccountsDatabase) GetAccountWithId(ctx context.Context, id int64) (*Account, error) {
	return nil, ErrNotFound
}

func (db *NullAccountsDatabase) GetAccountWithName(ctx context.Context, name string) (*Account, error) {
	return nil, ErrNotFound
}

func (db *NullAccountsDatabase) UpdateAccount(ctx context.Context, acct *Account) error {
	return nil
}

func (db *NullAccountsDatabase) RemoveAccount(ctx context.Context, acct *Account) error {
	return nil
}

func (db *NullAccountsDatabase) Close(ctx context.Context) error {
	return nil
}
