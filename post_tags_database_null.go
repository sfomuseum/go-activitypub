package activitypub

import (
	"context"
)

type NullPostTagsDatabase struct {
	PostTagsDatabase
}

func init() {
	ctx := context.Background()
	RegisterPostTagsDatabase(ctx, "null", NewNullPostTagsDatabase)
}

func NewNullPostTagsDatabase(ctx context.Context, uri string) (PostTagsDatabase, error) {
	db := &NullPostTagsDatabase{}
	return db, nil
}

func (db *NullPostTagsDatabase) GetLikeIdsForDateRange(ctx context.Context, start int64, end int64, cb GetPostTagIdsCallbackFunc) error {
	return nil
}

func (db *NullPostTagsDatabase) GetPostTagWithId(ctx context.Context, id int64) (*PostTag, error) {
	return nil, ErrNotFound
}

func (db *NullPostTagsDatabase) GetPostTagsForName(ctx context.Context, name string, cb GetPostTagsCallbackFunc) error {
	return nil
}

func (db *NullPostTagsDatabase) GetPostTagsForAccount(ctx context.Context, account_id int64, cb GetPostTagsCallbackFunc) error {
	return nil
}

func (db *NullPostTagsDatabase) GetPostTagsForPost(ctx context.Context, post_id int64, cb GetPostTagsCallbackFunc) error {
	return nil
}

func (db *NullPostTagsDatabase) AddPostTag(ctx context.Context, boost *PostTag) error {
	return nil
}

func (db *NullPostTagsDatabase) RemovePostTag(ctx context.Context, boost *PostTag) error {
	return nil
}

func (db *NullPostTagsDatabase) Close(ctx context.Context) error {
	return nil
}
