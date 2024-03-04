package activitypub

import (
	"context"
)

type NullLikesDatabase struct {
	LikesDatabase
}

func init() {
	ctx := context.Background()
	RegisterLikesDatabase(ctx, "null", NewNullLikesDatabase)
}

func NewNullLikesDatabase(ctx context.Context, uri string) (LikesDatabase, error) {
	db := &NullLikesDatabase{}
	return db, nil
}

func (db *NullLikesDatabase) GetLikeWithId(ctx context.Context, id int64) (*Like, error) {
	return nil, ErrNotFound
}

func (db *NullLikesDatabase) GetLikeWithPostIdAndActor(ctx context.Context, id int64, actor string) (*Like, error) {
	return nil, ErrNotFound
}

func (db *NullLikesDatabase) GetLikesForPost(ctx context.Context, post_id int64, cb GetLikesCallbackFunc) error {
	return nil
}

func (db *NullLikesDatabase) AddLike(ctx context.Context, like *Like) error {
	return nil
}

func (db *NullLikesDatabase) RemoveLike(ctx context.Context, like *Like) error {
	return nil
}

func (db *NullLikesDatabase) Close(ctx context.Context) error {
	return nil
}
