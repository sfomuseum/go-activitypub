package database

import (
	"context"
	"iter"

	"github.com/sfomuseum/go-activitypub"
)

type NullAccountsDatabase struct {
	Database[*activitypub.Account]
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

func (db *NullAccountsDatabase) AddRecord(ctx context.Context, a *activitypub.Account) error {
	return nil
}

func (db *NullAccountsDatabase) UpdateRecord(ctx context.Context, acct *activitypub.Account) error {
	return nil
}

func (db *NullAccountsDatabase) RemoveRecord(ctx context.Context, acct *activitypub.Account) error {
	return nil
}

func (db *NullAccountsDatabase) GetRecord(ctx context.Context, id int64) (*activitypub.Account, error) {
	return nil, activitypub.ErrNotFound
}

func (db *NullAccountsDatabase) QueryRecords(ctx context.Context, q *Query) iter.Seq2[*activitypub.Account, error] {
	return func(yield func(*activitypub.Account, error) bool) {}
}

func (db *NullAccountsDatabase) Close() error {
	return nil
}

func (db *NullAccountsDatabase) GetAccountIdsForDateRange(ctx context.Context, start int64, end int64) iter.Seq2[int64, error] {
	return func(yield func(int64, error) bool) {}
}

func (db *NullAccountsDatabase) GetAccountWithName(ctx context.Context, name string) (*activitypub.Account, error) {
	return nil, activitypub.ErrNotFound
}
