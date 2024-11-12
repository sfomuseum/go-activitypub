package database

import (
	"context"
	"fmt"
	"io"

	aa_docstore "github.com/aaronland/gocloud-docstore"
	"github.com/sfomuseum/go-activitypub"
	gc_docstore "gocloud.dev/docstore"
)

type DocstoreBoostsDatabase struct {
	BoostsDatabase
	collection *gc_docstore.Collection
}

func init() {

	ctx := context.Background()

	err := RegisterBoostsDatabase(ctx, "awsdynamodb", NewDocstoreBoostsDatabase)

	if err != nil {
		panic(err)
	}

	for _, scheme := range gc_docstore.DefaultURLMux().CollectionSchemes() {
		err := RegisterBoostsDatabase(ctx, scheme, NewDocstoreBoostsDatabase)

		if err != nil {
			panic(err)
		}
	}
}

func NewDocstoreBoostsDatabase(ctx context.Context, uri string) (BoostsDatabase, error) {

	col, err := aa_docstore.OpenCollection(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to open collection, %w", err)
	}

	db := &DocstoreBoostsDatabase{
		collection: col,
	}

	return db, nil
}

func (db *DocstoreBoostsDatabase) GetBoostIdsForDateRange(ctx context.Context, start int64, end int64, cb GetBoostIdsCallbackFunc) error {

	q := db.collection.Query()
	q = q.Where("Created", ">=", start)
	q = q.Where("Created", "<=", end)

	iter := q.Get(ctx, "Id")
	defer iter.Stop()

	for {

		var b activitypub.Boost
		err := iter.Next(ctx, &b)

		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("Failed to interate, %w", err)
		} else {
			err := cb(ctx, b.Id)

			if err != nil {
				return fmt.Errorf("Failed to invoke callback for boost %d, %w", b.Id, err)
			}
		}
	}

	return nil
}

func (db *DocstoreBoostsDatabase) GetBoostWithId(ctx context.Context, id int64) (*activitypub.Boost, error) {

	q := db.collection.Query()
	q = q.Where("Id", "=", id)

	return db.getBoost(ctx, q)
}

func (db *DocstoreBoostsDatabase) GetBoostWithPostIdAndActor(ctx context.Context, post_id int64, actor string) (*activitypub.Boost, error) {

	q := db.collection.Query()
	q = q.Where("PostId", "=", post_id)
	q = q.Where("Actor", "=", actor)

	return db.getBoost(ctx, q)
}

func (db *DocstoreBoostsDatabase) GetBoostsForPost(ctx context.Context, post_id int64, cb GetBoostsCallbackFunc) error {

	q := db.collection.Query()
	q = q.Where("PostId", "=", post_id)

	return db.getBoostsForQuery(ctx, q, cb)
}

func (db *DocstoreBoostsDatabase) GetBoostsForAccount(ctx context.Context, account_id int64, cb GetBoostsCallbackFunc) error {

	q := db.collection.Query()
	q = q.Where("AccountId", "=", account_id)

	return db.getBoostsForQuery(ctx, q, cb)
}

func (db *DocstoreBoostsDatabase) AddBoost(ctx context.Context, boost *activitypub.Boost) error {

	return db.collection.Put(ctx, boost)
}

func (db *DocstoreBoostsDatabase) RemoveBoost(ctx context.Context, boost *activitypub.Boost) error {

	return db.collection.Delete(ctx, boost)
}

func (db *DocstoreBoostsDatabase) Close(ctx context.Context) error {
	return db.collection.Close()
}

func (db *DocstoreBoostsDatabase) getBoost(ctx context.Context, q *gc_docstore.Query) (*activitypub.Boost, error) {

	iter := q.Get(ctx)
	defer iter.Stop()

	var b activitypub.Boost
	err := iter.Next(ctx, &b)

	if err == io.EOF {
		return nil, activitypub.ErrNotFound
	} else if err != nil {
		return nil, fmt.Errorf("Failed to interate, %w", err)
	} else {
		return &b, nil
	}

}

func (db *DocstoreBoostsDatabase) getBoostsForQuery(ctx context.Context, q *gc_docstore.Query, cb GetBoostsCallbackFunc) error {

	iter := q.Get(ctx)
	defer iter.Stop()

	for {

		var b activitypub.Boost
		err := iter.Next(ctx, &b)

		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("Failed to interate, %w", err)
		} else {

			err := cb(ctx, &b)

			if err != nil {
				return fmt.Errorf("Failed to execute callback for boost %d, %w", b.Id, err)
			}
		}
	}

	return nil
}
