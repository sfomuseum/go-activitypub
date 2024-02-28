package activitypub

import (
	"context"
	"fmt"
	"io"

	aa_docstore "github.com/aaronland/gocloud-docstore"
	gc_docstore "gocloud.dev/docstore"
)

type DocstoreAliasesDatabase struct {
	AliasesDatabase
	collection *gc_docstore.Collection
}

func init() {
	ctx := context.Background()

	RegisterAliasesDatabase(ctx, "awsdynamodb", NewDocstoreAliasesDatabase)

	for _, scheme := range gc_docstore.DefaultURLMux().CollectionSchemes() {
		RegisterAliasesDatabase(ctx, scheme, NewDocstoreAliasesDatabase)
	}
}

func NewDocstoreAliasesDatabase(ctx context.Context, uri string) (AliasesDatabase, error) {

	col, err := aa_docstore.OpenCollection(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to open collection, %w", err)
	}

	db := &DocstoreAliasesDatabase{
		collection: col,
	}

	return db, nil
}

func (db *DocstoreAliasesDatabase) GetAliasesForAccount(ctx context.Context, account_id int64, cb GetAliasesCallbackFunc) error {

	q := db.collection.Query()
	q = q.Where("AccountId", "=", account_id)

	iter := q.Get(ctx)
	defer iter.Stop()

	for {

		var a Alias
		err := iter.Next(ctx, &a)

		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("Failed to interate, %w", err)
		} else {

			err := cb(ctx, &a)

			if err != nil {
				return fmt.Errorf("Failed to execute callback for alias %s, %w", a, err)
			}
		}
	}

	return nil
}

func (db *DocstoreAliasesDatabase) GetAliasWithName(ctx context.Context, name string) (*Alias, error) {

	q := db.collection.Query()
	q = q.Where("Name", "=", name)

	return db.getAlias(ctx, q)
}

func (db *DocstoreAliasesDatabase) AddAlias(ctx context.Context, alias *Alias) error {
	return db.collection.Put(ctx, alias)
}

func (db *DocstoreAliasesDatabase) RemoveAlias(ctx context.Context, alias *Alias) error {
	return db.collection.Delete(ctx, alias)
}

func (db *DocstoreAliasesDatabase) Close(ctx context.Context) error {
	return db.collection.Close()
}

func (db *DocstoreAliasesDatabase) getAlias(ctx context.Context, q *gc_docstore.Query) (*Alias, error) {

	iter := q.Get(ctx)
	defer iter.Stop()

	var a Alias
	err := iter.Next(ctx, &a)

	if err == io.EOF {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, fmt.Errorf("Failed to interate, %w", err)
	} else {
		return &a, nil
	}
}
