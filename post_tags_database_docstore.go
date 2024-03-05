package activitypub

import (
	"context"
	"fmt"
	"io"

	aa_docstore "github.com/aaronland/gocloud-docstore"
	gc_docstore "gocloud.dev/docstore"
)

type DocstorePostTagsDatabase struct {
	PostTagsDatabase
	collection *gc_docstore.Collection
}

func init() {

	ctx := context.Background()

	RegisterPostTagsDatabase(ctx, "awsdynamodb", NewDocstorePostTagsDatabase)

	for _, scheme := range gc_docstore.DefaultURLMux().CollectionSchemes() {
		RegisterPostTagsDatabase(ctx, scheme, NewDocstorePostTagsDatabase)
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

func (db *DocstorePostTagsDatabase) GetPostTagWithId(ctx context.Context, id int64) (*PostTag, error) {
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

func (db *DocstorePostTagsDatabase) AddPostTag(ctx context.Context, tag *PostTag) error {
	return db.collection.Delete(ctx, tag)
}

func (db *DocstorePostTagsDatabase) RemovePostTag(ctx context.Context, boost *PostTag) error {
	return db.collection.Close()
}

func (db *DocstorePostTagsDatabase) Close(ctx context.Context) error {
	return nil
}

func (db *DocstorePostTagsDatabase) getPostTag(ctx context.Context, q *gc_docstore.Query) (*PostTag, error) {

	iter := q.Get(ctx)
	defer iter.Stop()

	var t PostTag
	err := iter.Next(ctx, &t)

	if err == io.EOF {
		return nil, ErrNotFound
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

		var t PostTag
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
