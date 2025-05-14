package database

import (
	"context"

	"github.com/sfomuseum/go-activitypub"
)

type NullPostsDatabase struct {
	PostsDatabase
}

func init() {

	ctx := context.Background()
	err := RegisterPostsDatabase(ctx, "null", NewNullPostsDatabase)

	if err != nil {
		panic(err)
	}
}

func NewNullPostsDatabase(ctx context.Context, uri string) (PostsDatabase, error) {

	db := &NullPostsDatabase{}
	return db, nil
}

func (db *NullPostsDatabase) GetPostIdsForDateRange(ctx context.Context, start int64, end int64, cb GetPostIdsCallbackFunc) error {
	return nil
}

func (db *NullPostsDatabase) AddPost(ctx context.Context, p *activitypub.Post) error {
	return nil
}

func (db *NullPostsDatabase) UpdatePost(ctx context.Context, p *activitypub.Post) error {
	return nil
}

func (db *NullPostsDatabase) GetPostWithId(ctx context.Context, id int64) (*activitypub.Post, error) {
	return nil, activitypub.ErrNotFound
}

func (db *NullPostsDatabase) Close(ctx context.Context) error {
	return nil
}
