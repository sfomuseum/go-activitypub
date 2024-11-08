package database

import (
	"context"
	"fmt"
	"io"

	aa_docstore "github.com/aaronland/gocloud-docstore"
	"github.com/sfomuseum/go-activitypub"
	gc_docstore "gocloud.dev/docstore"
)

type DocstoreLikesDatabase struct {
	LikesDatabase
	collection *gc_docstore.Collection
}

func init() {

	ctx := context.Background()

	err := RegisterLikesDatabase(ctx, "awsdynamodb", NewDocstoreLikesDatabase)

	if err != nil {
		panic(err)
	}

	for _, scheme := range gc_docstore.DefaultURLMux().CollectionSchemes() {
		err := RegisterLikesDatabase(ctx, scheme, NewDocstoreLikesDatabase)

		if err != nil {
			panic(err)
		}
	}
}

func NewDocstoreLikesDatabase(ctx context.Context, uri string) (LikesDatabase, error) {

	col, err := aa_docstore.OpenCollection(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to open collection, %w", err)
	}

	db := &DocstoreLikesDatabase{
		collection: col,
	}

	return db, nil
}

func (db *DocstoreLikesDatabase) GetLikeIdsForDateRange(ctx context.Context, start int64, end int64, cb GetLikeIdsCallbackFunc) error {

	q := db.collection.Query()
	q = q.Where("Created", ">=", start)
	q = q.Where("Created", "<=", end)

	iter := q.Get(ctx, "Id")
	defer iter.Stop()

	for {

		var l activitypub.Like
		err := iter.Next(ctx, &l)

		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("Failed to interate, %w", err)
		} else {
			err := cb(ctx, l.Id)

			if err != nil {
				return fmt.Errorf("Failed to invoke callback for like %d, %w", l.Id, err)
			}
		}
	}

	return nil
}

func (db *DocstoreLikesDatabase) GetLikeWithId(ctx context.Context, id int64) (*activitypub.Like, error) {

	q := db.collection.Query()
	q = q.Where("Id", "=", id)

	return db.getLike(ctx, q)
}

func (db *DocstoreLikesDatabase) GetLikeWithPostIdAndActor(ctx context.Context, post_id int64, actor string) (*activitypub.Like, error) {

	q := db.collection.Query()
	q = q.Where("PostId", "=", post_id)
	q = q.Where("Actor", "=", actor)

	return db.getLike(ctx, q)
}

func (db *DocstoreLikesDatabase) GetLikesForPost(ctx context.Context, post_id int64, cb GetLikesCallbackFunc) error {

	q := db.collection.Query()
	q = q.Where("PostId", "=", post_id)

	iter := q.Get(ctx)
	defer iter.Stop()

	for {

		var b activitypub.Like
		err := iter.Next(ctx, &b)

		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("Failed to interate, %w", err)
		} else {

			err := cb(ctx, &b)

			if err != nil {
				return fmt.Errorf("Failed to execute callback for like %d, %w", b.Id, err)
			}
		}
	}

	return nil

}

func (db *DocstoreLikesDatabase) getLike(ctx context.Context, q *gc_docstore.Query) (*activitypub.Like, error) {

	iter := q.Get(ctx)
	defer iter.Stop()

	var b activitypub.Like
	err := iter.Next(ctx, &b)

	if err == io.EOF {
		return nil, activitypub.ErrNotFound
	} else if err != nil {
		return nil, fmt.Errorf("Failed to interate, %w", err)
	} else {
		return &b, nil
	}

}

func (db *DocstoreLikesDatabase) AddLike(ctx context.Context, like *activitypub.Like) error {

	return db.collection.Put(ctx, like)
}

func (db *DocstoreLikesDatabase) RemoveLike(ctx context.Context, like *activitypub.Like) error {

	return db.collection.Delete(ctx, like)
}

func (db *DocstoreLikesDatabase) Close(ctx context.Context) error {
	return db.collection.Close()
}
