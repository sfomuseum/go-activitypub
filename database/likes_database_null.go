package database

import (
	"context"

	"github.com/sfomuseum/go-activitypub"
)

type NullLikesDatabase struct {
	LikesDatabase
}

func init() {
	ctx := context.Background()
	err := RegisterLikesDatabase(ctx, "null", NewNullLikesDatabase)

	if err != nil {
		panic(err)
	}
}

func NewNullLikesDatabase(ctx context.Context, uri string) (LikesDatabase, error) {
	db := &NullLikesDatabase{}
	return db, nil
}

func (db *NullLikesDatabase) GetLikeIdsForDateRange(ctx context.Context, start int64, end int64, cb GetLikeIdsCallbackFunc) error {
	return nil
}

func (db *NullLikesDatabase) GetLikeWithId(ctx context.Context, id int64) (*activitypub.Like, error) {
	return nil, activitypub.ErrNotFound
}

func (db *NullLikesDatabase) GetLikeWithPostIdAndActor(ctx context.Context, id int64, actor string) (*activitypub.Like, error) {
	return nil, activitypub.ErrNotFound
}

func (db *NullLikesDatabase) GetLikesForPost(ctx context.Context, post_id int64, cb GetLikesCallbackFunc) error {
	return nil
}

func (db *NullLikesDatabase) AddLike(ctx context.Context, like *activitypub.Like) error {
	return nil
}

func (db *NullLikesDatabase) RemoveLike(ctx context.Context, like *activitypub.Like) error {
	return nil
}

func (db *NullLikesDatabase) Close(ctx context.Context) error {
	return nil
}
