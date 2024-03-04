package activitypub

import (
	"context"
	"fmt"
	"io"

	aa_docstore "github.com/aaronland/gocloud-docstore"
	gc_docstore "gocloud.dev/docstore"
)

type DocstoreBoostsDatabase struct {
	BoostsDatabase
	collection *gc_docstore.Collection
}

func init() {

	ctx := context.Background()

	RegisterBoostsDatabase(ctx, "awsdynamodb", NewDocstoreBoostsDatabase)

	for _, scheme := range gc_docstore.DefaultURLMux().CollectionSchemes() {
		RegisterBoostsDatabase(ctx, scheme, NewDocstoreBoostsDatabase)
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

func (db *DocstoreBoostsDatabase) GetBoostWithId(ctx context.Context, id int64) (*Boost, error) {

	q := db.collection.Query()
	q = q.Where("Id", "=", id)

	return db.getBoost(ctx, q)
}

func (db *DocstoreBoostsDatabase) GetBoostWithPostIdAndActor(ctx context.Context, post_id int64, actor string) (*Boost, error) {

	q := db.collection.Query()
	q = q.Where("PostId", "=", post_id)
	q = q.Where("Actor", "=", actor)

	return db.getBoost(ctx, q)
}

func (db *DocstoreBoostsDatabase) GetBoostsForPost(ctx context.Context, post_id int64, cb GetBoostsCallbackFunc) error {

	q := db.collection.Query()
	q = q.Where("PostId", "=", post_id)

	iter := q.Get(ctx)
	defer iter.Stop()

	for {

		var b Boost
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

func (db *DocstoreBoostsDatabase) getBoost(ctx context.Context, q *gc_docstore.Query) (*Boost, error) {

	iter := q.Get(ctx)
	defer iter.Stop()

	var b Boost
	err := iter.Next(ctx, &b)

	if err == io.EOF {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, fmt.Errorf("Failed to interate, %w", err)
	} else {
		return &b, nil
	}

}

func (db *DocstoreBoostsDatabase) AddBoost(ctx context.Context, boost *Boost) error {

	return db.collection.Put(ctx, boost)
}

func (db *DocstoreBoostsDatabase) RemoveBoost(ctx context.Context, boost *Boost) error {

	return db.collection.Delete(ctx, boost)
}

func (db *DocstoreBoostsDatabase) Close(ctx context.Context) error {
	return db.collection.Close()
}
