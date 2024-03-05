package activitypub

import (
	"context"
)

type NullRepliesDatabase struct {
	RepliesDatabase
}

func init() {
	ctx := context.Background()
	RegisterRepliesDatabase(ctx, "null", NewNullRepliesDatabase)
}

func NewNullRepliesDatabase(ctx context.Context, uri string) (RepliesDatabase, error) {
	db := &NullRepliesDatabase{}
	return db, nil
}

func (db *NullRepliesDatabase) GetReplyWithId(ctx context.Context, id int64) (*Reply, error) {
	return nil, ErrNotFound
}

func (db *NullRepliesDatabase) GetReplyWithReplyId(ctx context.Context, reply_id string) (*Reply, error) {
	return nil, ErrNotFound
}

func (db *NullRepliesDatabase) GetRepliesForPost(ctx context.Context, post_id int64, cb GetRepliesCallbackFunc) error {
	return nil
}

func (db *NullRepliesDatabase) GetRepliesForAccount(ctx context.Context, post_id int64, cb GetRepliesCallbackFunc) error {
	return nil
}

func (db *NullRepliesDatabase) GetRepliesForActor(ctx context.Context, actor string, cb GetRepliesCallbackFunc) error {
	return nil
}

func (db *NullRepliesDatabase) AddReply(ctx context.Context, reply *Reply) error {
	return nil
}

func (db *NullRepliesDatabase) RemoveReply(ctx context.Context, reply *Reply) error {
	return nil
}

func (db *NullRepliesDatabase) Close(ctx context.Context) error {
	return nil
}
