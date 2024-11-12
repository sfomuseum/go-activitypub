package database

import (
	"context"
	"fmt"
	"io"

	aa_docstore "github.com/aaronland/gocloud-docstore"
	"github.com/sfomuseum/go-activitypub"
	gc_docstore "gocloud.dev/docstore"
)

type DocstorePostTagsDatabase struct {
	PostTagsDatabase
	collection *gc_docstore.Collection
}

func init() {

	ctx := context.Background()

	err := RegisterPostTagsDatabase(ctx, "awsdynamodb", NewDocstorePostTagsDatabase)

	if err != nil {
		panic(err)
	}

	for _, scheme := range gc_docstore.DefaultURLMux().CollectionSchemes() {
		err := RegisterPostTagsDatabase(ctx, scheme, NewDocstorePostTagsDatabase)

		if err != nil {
			panic(err)
		}
	}
}

func NewDocstorePostTagsDatabase(ctx context.Context, uri string) (PostTagsDatabase, error) {

	col, err := aa_docstore.OpenCollection(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to open collection, %w", err)
	}

	db := &DocstorePostTagsDatabase{
		collection: col,
	}
	return db, nil
}

func (db *DocstorePostTagsDatabase) GetPostTagIdsForDateRange(ctx context.Context, start int64, end int64, cb GetPostTagIdsCallbackFunc) error {

	q := db.collection.Query()
	q = q.Where("Created", ">=", start)
	q = q.Where("Created", "<=", end)

	iter := q.Get(ctx, "Id")
	defer iter.Stop()

	for {

		var t activitypub.PostTag
		err := iter.Next(ctx, &t)

		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("Failed to interate, %w", err)
		} else {
			err := cb(ctx, t.Id)

			if err != nil {
				return fmt.Errorf("Failed to invoke callback for post tag %d, %w", t.Id, err)
			}
		}
	}

	return nil
}

func (db *DocstorePostTagsDatabase) GetPostTagWithId(ctx context.Context, id int64) (*activitypub.PostTag, error) {
	q := db.collection.Query()
	q = q.Where("Id", "=", id)
	return db.getPostTag(ctx, q)
}

func (db *DocstorePostTagsDatabase) GetPostTagsForName(ctx context.Context, name string, cb GetPostTagsCallbackFunc) error {
	q := db.collection.Query()
	q = q.Where("Name", "=", name)
	return db.getPostTags(ctx, q, cb)
}

func (db *DocstorePostTagsDatabase) GetPostTagsForAccount(ctx context.Context, account_id int64, cb GetPostTagsCallbackFunc) error {
	q := db.collection.Query()
	q = q.Where("AccountId", "=", account_id)
	return db.getPostTags(ctx, q, cb)
}

func (db *DocstorePostTagsDatabase) GetPostTagsForPost(ctx context.Context, post_id int64, cb GetPostTagsCallbackFunc) error {
	q := db.collection.Query()
	q = q.Where("PostId", "=", post_id)
	return db.getPostTags(ctx, q, cb)
}

func (db *DocstorePostTagsDatabase) AddPostTag(ctx context.Context, tag *activitypub.PostTag) error {
	return db.collection.Put(ctx, tag)
}

func (db *DocstorePostTagsDatabase) RemovePostTag(ctx context.Context, tag *activitypub.PostTag) error {
	return db.collection.Delete(ctx, tag)
}

func (db *DocstorePostTagsDatabase) Close(ctx context.Context) error {
	return db.collection.Close()
}

func (db *DocstorePostTagsDatabase) getPostTag(ctx context.Context, q *gc_docstore.Query) (*activitypub.PostTag, error) {

	iter := q.Get(ctx)
	defer iter.Stop()

	var t activitypub.PostTag
	err := iter.Next(ctx, &t)

	if err == io.EOF {
		return nil, activitypub.ErrNotFound
	} else if err != nil {
		return nil, fmt.Errorf("Failed to interate, %w", err)
	} else {
		return &t, nil
	}
}

func (db *DocstorePostTagsDatabase) getPostTags(ctx context.Context, q *gc_docstore.Query, cb GetPostTagsCallbackFunc) error {

	iter := q.Get(ctx)
	defer iter.Stop()

	for {

		var t activitypub.PostTag
		err := iter.Next(ctx, &t)

		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("Failed to interate, %w", err)
		} else {

			err := cb(ctx, &t)

			if err != nil {
				return fmt.Errorf("Failed to execute callback for tag %d, %w", t.Id, err)
			}
		}
	}

	return nil
}
