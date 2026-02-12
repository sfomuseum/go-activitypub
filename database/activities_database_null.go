package database

import (
	"context"
	"iter"

	"github.com/sfomuseum/go-activitypub"
)

type NullActivitiesDatabase struct {
	Database[*activitypub.Activity]
	ActivitiesDatabase
}

func init() {
	ctx := context.Background()
	err := RegisterActivitiesDatabase(ctx, "null", NewNullActivitiesDatabase)

	if err != nil {
		panic(err)
	}
}

func NewNullActivitiesDatabase(ctx context.Context, uri string) (ActivitiesDatabase, error) {
	db := &NullActivitiesDatabase{}
	return db, nil
}

func (db *NullActivitiesDatabase) AddRecord(ctx context.Context, f *activitypub.Activity) error {
	return nil
}

func (db *NullActivitiesDatabase) UpdateRecord(ctx context.Context, f *activitypub.Activity) error {
	return nil
}

func (db *NullActivitiesDatabase) RemoveRecord(ctx context.Context, f *activitypub.Activity) error {
	return nil
}

func (db *NullActivitiesDatabase) GetRecord(ctx context.Context, id int64) (*activitypub.Activity, error) {
	return nil, activitypub.ErrNotFound
}

func (db *NullActivitiesDatabase) QueryRecords(ctx context.Context, q *Query) iter.Seq2[*activitypub.Activity, error] {
	return func(yield func(*activitypub.Activity, error) bool) {}
}

func (db *NullActivitiesDatabase) Close() error {
	return nil
}

func (db *NullActivitiesDatabase) GetActivityWithActivityPubId(ctx context.Context, id string) (*activitypub.Activity, error) {
	return nil, activitypub.ErrNotFound
}

func (db *NullActivitiesDatabase) GetActivityWithActivityTypeAndId(ctx context.Context, activity_type activitypub.ActivityType, id int64) (*activitypub.Activity, error) {
	return nil, activitypub.ErrNotFound
}

func (db *NullActivitiesDatabase) GetActivitiesForAccount(ctx context.Context, id int64) iter.Seq2[*activitypub.Activity, error] {
	return func(yield func(*activitypub.Activity, error) bool) {}
}
