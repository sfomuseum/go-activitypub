package activitypub

import (
	"context"
	"fmt"
	"io"

	aa_docstore "github.com/aaronland/gocloud-docstore"
	gc_docstore "gocloud.dev/docstore"
)

type DocstoreRepliesDatabase struct {
	RepliesDatabase
	collection *gc_docstore.Collection
}

func init() {

	ctx := context.Background()

	RegisterRepliesDatabase(ctx, "awsdynamodb", NewDocstoreRepliesDatabase)

	for _, scheme := range gc_docstore.DefaultURLMux().CollectionSchemes() {
		RegisterRepliesDatabase(ctx, scheme, NewDocstoreRepliesDatabase)
	}
}

func NewDocstoreRepliesDatabase(ctx context.Context, uri string) (RepliesDatabase, error) {

	col, err := aa_docstore.OpenCollection(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to open collection, %w", err)
	}

	db := &DocstoreRepliesDatabase{
		collection: col,
	}
	return db, nil
}

func (db *DocstoreRepliesDatabase) GetReplyWithId(ctx context.Context, id int64) (*Reply, error) {
	q := db.collection.Query()
	q = q.Where("Id", "=", id)
	return db.getReply(ctx, q)
}

func (db *DocstoreRepliesDatabase) GetReplyWithReplyId(ctx context.Context, reply_id string) (*Reply, error) {
	q := db.collection.Query()
	q = q.Where("ReplyId", "=", reply_id)
	return db.getReply(ctx, q)
}

func (db *DocstoreRepliesDatabase) GetRepliesForPost(ctx context.Context, post_id int64, cb GetRepliesCallbackFunc) error {
	q := db.collection.Query()
	q = q.Where("PostId", "=", post_id)
	return db.getReplies(ctx, q, cb)
}

func (db *DocstoreRepliesDatabase) GetRepliesForAccount(ctx context.Context, account_id int64, cb GetRepliesCallbackFunc) error {
	q := db.collection.Query()
	q = q.Where("AccountId", "=", account_id)
	return db.getReplies(ctx, q, cb)
}

func (db *DocstoreRepliesDatabase) GetRepliesForActor(ctx context.Context, actor string, cb GetRepliesCallbackFunc) error {
	q := db.collection.Query()
	q = q.Where("Actor", "=", actor)
	return db.getReplies(ctx, q, cb)
}

func (db *DocstoreRepliesDatabase) AddReply(ctx context.Context, reply *Reply) error {
	return db.collection.Put(ctx, reply)
}

func (db *DocstoreRepliesDatabase) RemoveReply(ctx context.Context, reply *Reply) error {
	return db.collection.Delete(ctx, reply)
}

func (db *DocstoreRepliesDatabase) Close(ctx context.Context) error {
	return db.collection.Close()
}

func (db *DocstoreRepliesDatabase) getReply(ctx context.Context, q *gc_docstore.Query) (*Reply, error) {

	iter := q.Get(ctx)
	defer iter.Stop()

	var r Reply
	err := iter.Next(ctx, &r)

	if err == io.EOF {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, fmt.Errorf("Failed to interate, %w", err)
	} else {
		return &r, nil
	}
}

func (db *DocstoreRepliesDatabase) getReplies(ctx context.Context, q *gc_docstore.Query, cb GetRepliesCallbackFunc) error {

	iter := q.Get(ctx)
	defer iter.Stop()

	for {

		var r Reply
		err := iter.Next(ctx, &r)

		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("Failed to interate, %w", err)
		} else {

			err := cb(ctx, &r)

			if err != nil {
				return fmt.Errorf("Failed to execute callback for reply %d, %w", r.Id, err)
			}
		}
	}

	return nil
}
