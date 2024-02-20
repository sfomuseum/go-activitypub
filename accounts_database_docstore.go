package activitypub

import (
	"context"
	"fmt"
	"io"

	aa_docstore "github.com/aaronland/gocloud-docstore"
	gc_docstore "gocloud.dev/docstore"
)

type DocstoreAccountsDatabase struct {
	AccountsDatabase
	collection *gc_docstore.Collection
}

func init() {

	ctx := context.Background()

	RegisterAccountsDatabase(ctx, "awsdynamodb", NewDocstoreAccountsDatabase)

	for _, scheme := range gc_docstore.DefaultURLMux().CollectionSchemes() {
		RegisterAccountsDatabase(ctx, scheme, NewDocstoreAccountsDatabase)
	}
}

func NewDocstoreAccountsDatabase(ctx context.Context, uri string) (AccountsDatabase, error) {

	col, err := aa_docstore.OpenCollection(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to open collection, %w", err)
	}

	db := &DocstoreAccountsDatabase{
		collection: col,
	}

	return db, nil
}

func (db *DocstoreAccountsDatabase) AddAccount(ctx context.Context, a *Account) error {

	return db.collection.Put(ctx, a)
}

func (db *DocstoreAccountsDatabase) GetAccountWithId(ctx context.Context, id int64) (*Account, error) {

	q := db.collection.Query()
	q = q.Where("Id", "=", id)

	return db.getAccount(ctx, q)
}

func (db *DocstoreAccountsDatabase) GetAccountWithName(ctx context.Context, name string) (*Account, error) {

	q := db.collection.Query()
	q = q.Where("Name", "=", name)

	return db.getAccount(ctx, q)
}

func (db *DocstoreAccountsDatabase) getAccount(ctx context.Context, q *gc_docstore.Query) (*Account, error) {

	iter := q.Get(ctx)
	defer iter.Stop()

	for {

		var a Account
		err := iter.Next(ctx, &a)

		if err == io.EOF {
			return nil, ErrNotFound
		} else if err != nil {
			return nil, fmt.Errorf("Failed to interate, %w", err)
		} else {
			return &a, nil
		}
	}

	return nil, ErrNotFound
}

func (db *DocstoreAccountsDatabase) Close(ctx context.Context) error {
	return db.collection.Close()
}
