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

func (db *DocstoreAccountsDatabase) GetAccountIdsForDateRange(ctx context.Context, start int64, end int64, cb GetAccountIdsCallbackFunc) error {

	q := db.collection.Query()

	// For reasons I don't understand this frequently panics along the lines of:

	/*
		panic: runtime error: index out of range [0] with length 0

		goroutine 98 [running]:
		gocloud.dev/docstore/awsdynamodb.(*documentIterator).Next(0x14000056050, {0x101ba36b0?, 0x10258bb40?}, {{0x101b386e0, 0x1400019bea0}, 0x0, {0x101b454e0, 0x1400019bea0, 0x199}, {0x140004fc808, ...}})
			/usr/local/sfomuseum/go-activitypub/vendor/gocloud.dev/docstore/awsdynamodb/query.go:492 +0x260
	*/

	// What I find most confusing is that this doesn't happen in any of the other _docstore
	// packages, for example deliveries_database_docstore.go

	// q = q.Where("Created", ">=", start)
	// q = q.Where("Created", "<=", end)

	// See also: https://github.com/google/go-cloud/issues/3405

	iter := q.Get(ctx)
	defer iter.Stop()

	for {

		var a Account
		err := iter.Next(ctx, &a)

		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("Failed to interate, %w", err)
		} else {

			if a.Created >= start && a.Created <= end { // START OF see notes wrt/panics above
				err := cb(ctx, a.Id)

				if err != nil {
					return fmt.Errorf("Failed to invoke callback for account %d, %w", a.Id, err)
				}
			} // END OF see notes wrt/panics above
		}
	}

	return nil
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

func (db *DocstoreAccountsDatabase) Close(ctx context.Context) error {
	return db.collection.Close()
}

func (db *DocstoreAccountsDatabase) getAccount(ctx context.Context, q *gc_docstore.Query) (*Account, error) {

	iter := q.Get(ctx)
	defer iter.Stop()

	var a Account
	err := iter.Next(ctx, &a)

	if err == io.EOF {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, fmt.Errorf("Failed to interate, %w", err)
	} else {
		return &a, nil
	}

	return nil, ErrNotFound
}
