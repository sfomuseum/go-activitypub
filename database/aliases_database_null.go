package database

import (
	"context"

	"github.com/sfomuseum/go-activitypub"
)

type NullAliasesDatabase struct {
	AliasesDatabase
}

func init() {
	ctx := context.Background()
	err := RegisterAliasesDatabase(ctx, "null", NewNullAliasesDatabase)

	if err != nil {
		panic(err)
	}

}

func NewNullAliasesDatabase(ctx context.Context, uri string) (AliasesDatabase, error) {
	db := &NullAliasesDatabase{}
	return db, nil
}

func (db *NullAliasesDatabase) GetAliasesForAccount(ctx context.Context, account_id int64, cb GetAliasesCallbackFunc) error {
	return nil
}

func (db *NullAliasesDatabase) GetAliasWithName(ctx context.Context, name string) (*activitypub.Alias, error) {
	return nil, activitypub.ErrNotFound
}

func (db *NullAliasesDatabase) AddAlias(ctx context.Context, alias *activitypub.Alias) error {
	return nil
}

func (db *NullAliasesDatabase) RemoveAlias(ctx context.Context, alias *activitypub.Alias) error {
	return nil
}

func (db *NullAliasesDatabase) Close(ctx context.Context) error {
	return nil
}
