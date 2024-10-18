package database

import (
	"context"
	"fmt"
	"io"

	aa_docstore "github.com/aaronland/gocloud-docstore"
	"github.com/sfomuseum/go-activitypub"
	gc_docstore "gocloud.dev/docstore"
)

type DocstoreActivitiesDatabase struct {
	ActivitiesDatabase
	collection *gc_docstore.Collection
}

func init() {

	ctx := context.Background()

	RegisterActivitiesDatabase(ctx, "awsdynamodb", NewDocstoreActivitiesDatabase)

	for _, scheme := range gc_docstore.DefaultURLMux().CollectionSchemes() {
		RegisterActivitiesDatabase(ctx, scheme, NewDocstoreActivitiesDatabase)
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

func (db *DocstoreActivitiesDatabase) AddActivity(ctx context.Context, f *activitypub.Activity) error {

	return db.collection.Put(ctx, f)
}

func (db *DocstoreActivitiesDatabase) GetActivityWithId(ctx context.Context, id int64) (*activitypub.Activity, error) {

	q := db.collection.Query()
	q = q.Where("Id", "=", id)

	return db.getActivity(ctx, q)
}

func (db *DocstoreActivitiesDatabase) Close(ctx context.Context) error {
	return db.collection.Close()
}

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
