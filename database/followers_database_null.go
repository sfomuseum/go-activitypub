package database

import (
	"context"

	"github.com/sfomuseum/go-activitypub"
)

type NullFollowersDatabase struct {
	FollowersDatabase
}

func init() {

	ctx := context.Background()
	err := RegisterFollowersDatabase(ctx, "null", NewNullFollowersDatabase)

	if err != nil {
		panic(err)
	}
}

func NewNullFollowersDatabase(ctx context.Context, uri string) (FollowersDatabase, error) {
	db := &NullFollowersDatabase{}
	return db, nil
}

func (db *NullFollowersDatabase) GetFollowerIdsForDateRange(ctx context.Context, start int64, end int64, cb GetFollowerIdsCallbackFunc) error {
	return nil
}

func (db *NullFollowersDatabase) GetAllFollowers(ctx context.Context, cb GetFollowersCallbackFunc) error {
	return nil
}

func (db *NullFollowersDatabase) HasFollowers(ctx context.Context, account_id int64) (bool, error) {
	return false, nil
}

func (db *NullFollowersDatabase) GetFollower(ctx context.Context, account_id int64, follower_address string) (*activitypub.Follower, error) {
	return nil, activitypub.ErrNotFound
}

func (db *NullFollowersDatabase) GetFollowerWithId(ctx context.Context, follower_id int64) (*activitypub.Follower, error) {
	return nil, activitypub.ErrNotFound
}

func (db *NullFollowersDatabase) Close(ctx context.Context) error {
	return nil
}
