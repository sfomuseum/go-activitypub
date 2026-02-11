package database

import (
	"context"
	"fmt"
	"io"
	"iter"
	
	aa_docstore "github.com/aaronland/gocloud/docstore"
	"github.com/sfomuseum/go-activitypub"
	gc_docstore "gocloud.dev/docstore"
)

type DocstoreAccountsDatabase struct {
	Database[*activitypub.Account]
	AccountsDatabase
	collection *gc_docstore.Collection
}

func init() {

	ctx := context.Background()

	err := RegisterAccountsDatabase(ctx, "awsdynamodb", NewDocstoreAccountsDatabase)

	if err != nil {
		panic(err)
	}

	for _, scheme := range gc_docstore.DefaultURLMux().CollectionSchemes() {
		err := RegisterAccountsDatabase(ctx, scheme, NewDocstoreAccountsDatabase)

		if err != nil {
			panic(err)
		}
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

func (db *DocstoreAccountsDatabase) Close() error {
	return db.collection.Close()
}

func (db *DocstoreAccountsDatabase) AddRecord(ctx context.Context, a *activitypub.Account) error {
	return db.collection.Put(ctx, a)
}

func (db *DocstoreAccountsDatabase) UpdateRecord(ctx context.Context, acct *activitypub.Account) error {
	return db.collection.Replace(ctx, acct)
}

func (db *DocstoreAccountsDatabase) RemoveRecord(ctx context.Context, acct *activitypub.Account) error {
	return db.collection.Delete(ctx, acct)
}

func (db *DocstoreAccountsDatabase) GetRecord(ctx context.Context, id int64) (*activitypub.Account, error) {
	q := db.collection.Query()
	q = q.Where("Id", "=", id)
	return db.getAccount(ctx, q)
}

func (db DocstoreAccountsDatabase) QueryRecords(ctx context.Context, q *Query) iter.Seq2[*activitypub.Account, error] {

	return func(yield func(*activitypub.Account, error) bool) {
		
		col_q := newDocstoreQuery(db.collection, q)
		
		iter := col_q.Get(ctx)
		defer iter.Stop()

		for {
			
			var a activitypub.Account
			err := iter.Next(ctx, &a)
			
			if err == io.EOF {
				break
			} else if err != nil {

				if !yield(nil, err){
					return
				}
				
			} else {

				if !yield(&a, nil){
					return
				}
			}
		}
	}
}

func (db *DocstoreAccountsDatabase) GetAccountIdsForDateRange(ctx context.Context, start int64, end int64) iter.Seq2[*activitypub.Account, error] {

	conditions := []*Condition{
		&Condition{
			Field: "Created",
			Operator: ">=",
			Value: start,
		},
		&Condition{
			Field: "Created",
			Operator: "<=",
			Value: end,
		},
	}

	where := &Where{
		Conditions: conditions,
		Relation: "AND",
	}
	
	q := &Query{
		Where: where,
	}

	return db.QueryRecords(ctx, q)
}

func (db *DocstoreAccountsDatabase) GetAccountWithName(ctx context.Context, name string) (*activitypub.Account, error) {

	q := db.collection.Query()
	q = q.Where("Name", "=", name)

	return db.getAccount(ctx, q)
}


func (db *DocstoreAccountsDatabase) getAccount(ctx context.Context, q *gc_docstore.Query) (*activitypub.Account, error) {

	iter := q.Get(ctx)
	defer iter.Stop()

	var a activitypub.Account
	err := iter.Next(ctx, &a)

	if err == io.EOF {
		return nil, activitypub.ErrNotFound
	} else if err != nil {
		return nil, fmt.Errorf("Failed to interate, %w", err)
	} else {
		return &a, nil
	}

	return nil, activitypub.ErrNotFound
}

