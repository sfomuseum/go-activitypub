package database

import (
	"context"

	"github.com/sfomuseum/go-activitypub"
)

type NullBoostsDatabase struct {
	BoostsDatabase
}

func init() {
	ctx := context.Background()
	err := RegisterBoostsDatabase(ctx, "null", NewNullBoostsDatabase)

	if err != nil {
		panic(err)
	}
}

func NewNullBoostsDatabase(ctx context.Context, uri string) (BoostsDatabase, error) {
	db := &NullBoostsDatabase{}
	return db, nil
}

func (db *NullBoostsDatabase) GetBoostIdsForDateRange(ctx context.Context, start int64, end int64, cb GetBoostIdsCallbackFunc) error {
	return nil
}

func (db *NullBoostsDatabase) GetBoostWithId(ctx context.Context, id int64) (*activitypub.Boost, error) {
	return nil, activitypub.ErrNotFound
}

func (db *NullBoostsDatabase) GetBoostWithPostIdAndActor(ctx context.Context, id int64, actor string) (*activitypub.Boost, error) {
	return nil, activitypub.ErrNotFound
}

func (db *NullBoostsDatabase) GetBoostsForPost(ctx context.Context, post_id int64, cb GetBoostsCallbackFunc) error {
	return nil
}

func (db *NullBoostsDatabase) GetBoostsForAccount(ctx context.Context, account_id int64, cb GetBoostsCallbackFunc) error {
	return nil
}

func (db *NullBoostsDatabase) AddBoost(ctx context.Context, boost *activitypub.Boost) error {
	return nil
}

func (db *NullBoostsDatabase) RemoveBoost(ctx context.Context, boost *activitypub.Boost) error {
	return nil
}

func (db *NullBoostsDatabase) Close(ctx context.Context) error {
	return nil
}
