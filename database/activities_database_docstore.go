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

type DocstoreActivitiesDatabase struct {
	Database[*activitypub.Activity]
	ActivitiesDatabase
	collection *gc_docstore.Collection
}

func init() {

	ctx := context.Background()

	err := RegisterActivitiesDatabase(ctx, "awsdynamodb", NewDocstoreActivitiesDatabase)

	if err != nil {
		panic(err)
	}

	for _, scheme := range gc_docstore.DefaultURLMux().CollectionSchemes() {
		err := RegisterActivitiesDatabase(ctx, scheme, NewDocstoreActivitiesDatabase)

		if err != nil {
			panic(err)
		}

	}
}

func NewDocstoreActivitiesDatabase(ctx context.Context, uri string) (ActivitiesDatabase, error) {

	col, err := aa_docstore.OpenCollection(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to open collection, %w", err)
	}

	db := &DocstoreActivitiesDatabase{
		collection: col,
	}

	return db, nil
}

// Database interface

func (db *DocstoreActivitiesDatabase) AddRecord(ctx context.Context, f *activitypub.Activity) error {
	return db.collection.Put(ctx, f)
}

func (db *DocstoreActivitiesDatabase) GetRecord(ctx context.Context, id int64) (*activitypub.Activity, error) {

	q := db.collection.Query()
	q = q.Where("Id", "=", id)

	return db.getActivity(ctx, q)
}

func (db *DocstoreActivitiesDatabase) UpdateRecord(ctx context.Context, a *activitypub.Activity) error {
	return activitypub.ErrNotImplemented
}

func (db *DocstoreActivitiesDatabase) RemoveRecord(ctx context.Context, a *activitypub.Activity) error {
	return activitypub.ErrNotImplemented
}

func (db *DocstoreActivitiesDatabase) QueryRecords(ctx context.Context, q *Query) iter.Seq2[*activitypub.Activity, error] {

	col_q := newDocstoreQuery(db.collection, q)
	return db.getActivitiesWithQuery(ctx, col_q)
}

func (db *DocstoreActivitiesDatabase) Close() error {
	return db.collection.Close()
}

// ActivitiesDatabase interface

func (db *DocstoreActivitiesDatabase) GetActivityWithActivityPubId(ctx context.Context, id string) (*activitypub.Activity, error) {

	q := db.collection.Query()
	q = q.Where("ActivityPubId", "=", id)

	return db.getActivity(ctx, q)
}

func (db *DocstoreActivitiesDatabase) GetActivityWithActivityTypeAndId(ctx context.Context, activity_type activitypub.ActivityType, id int64) (*activitypub.Activity, error) {

	q := db.collection.Query()
	q = q.Where("ActivityType", "=", activity_type)
	q = q.Where("ActivityTypeId", "=", id)

	return db.getActivity(ctx, q)
}

func (db *DocstoreActivitiesDatabase) GetActivitiesForAccount(ctx context.Context, id int64) iter.Seq2[*activitypub.Activity, error] {

	q := db.collection.Query()
	q = q.Where("AccountId", "=", id)

	return db.getActivitiesWithQuery(ctx, q)
}

// Local methods

func (db *DocstoreActivitiesDatabase) getActivity(ctx context.Context, q *gc_docstore.Query) (*activitypub.Activity, error) {

	iter := q.Get(ctx)
	defer iter.Stop()

	for {

		var d activitypub.Activity
		err := iter.Next(ctx, &d)

		if err == io.EOF {
			return nil, activitypub.ErrNotFound
		} else if err != nil {
			return nil, fmt.Errorf("Failed to interate, %w", err)
		} else {
			return &d, nil
		}
	}

	return nil, activitypub.ErrNotFound

}

func (db *DocstoreActivitiesDatabase) getActivitiesWithQuery(ctx context.Context, q *gc_docstore.Query) iter.Seq2[*activitypub.Activity, error] {

	return func(yield func(*activitypub.Activity, error) bool) {

		iter := q.Get(ctx)
		defer iter.Stop()

		for {

			var a activitypub.Activity
			err := iter.Next(ctx, &a)

			if err == io.EOF {
				break
			} else if err != nil {
				if !yield(nil, fmt.Errorf("Failed to interate, %w", err)) {
					return
				}
			} else {

				if !yield(&a, nil) {
					return
				}
			}
		}
	}

}
