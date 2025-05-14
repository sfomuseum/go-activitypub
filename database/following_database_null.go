package database

import (
	"context"
	"github.com/sfomuseum/go-activitypub"
)

type NullFollowingDatabase struct {
	FollowingDatabase
}

func init() {

	ctx := context.Background()
	err := RegisterFollowingDatabase(ctx, "null", NewNullFollowingDatabase)

	if err != nil {
		panic(err)
	}
}

func NewNullFollowingDatabase(ctx context.Context, uri string) (FollowingDatabase, error) {
	db := &NullFollowingDatabase{}
	return db, nil
}

func (db *NullFollowingDatabase) GetFollowingIdsForDateRange(ctx context.Context, start int64, end int64, cb GetFollowingIdsCallbackFunc) error {
	return nil
}

func (db *NullFollowingDatabase) GetFollowing(ctx context.Context, account_id int64, following_address string) (*activitypub.Following, error) {
	return nil, activitypub.ErrNotFound
}

func (db *NullFollowingDatabase) AddFollowing(ctx context.Context, f *activitypub.Following) error {
	return nil
}

func (db *NullFollowingDatabase) RemoveFollowing(ctx context.Context, f *activitypub.Following) error {
	return nil
}

func (db *NullFollowingDatabase) GetFollowingForAccount(ctx context.Context, account_id int64, following_callback GetFollowingCallbackFunc) error {
	return nil
}

func (db *NullFollowingDatabase) Close(ctx context.Context) error {
	return nil
}
