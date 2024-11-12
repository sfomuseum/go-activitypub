package database

import (
	"context"

	"github.com/sfomuseum/go-activitypub"
)

type NullActivitiesDatabase struct {
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

func (db *NullActivitiesDatabase) AddActivity(ctx context.Context, f *activitypub.Activity) error {
	return nil
}

func (db *NullActivitiesDatabase) GetActivityWithId(ctx context.Context, id int64) (*activitypub.Activity, error) {
	return nil, activitypub.ErrNotFound
}

func (db *NullActivitiesDatabase) GetActivityWithActivityPubId(ctx context.Context, id string) (*activitypub.Activity, error) {
	return nil, activitypub.ErrNotFound
}

func (db *NullActivitiesDatabase) GetActivityWithActivityTypeAndId(ctx context.Context, activity_type activitypub.ActivityType, id int64) (*activitypub.Activity, error) {
	return nil, activitypub.ErrNotFound
}

func (db *NullActivitiesDatabase) GetActivities(ctx context.Context, cb GetActivitiesCallbackFunc) error {
	return nil
}

func (db *NullActivitiesDatabase) GetActivitiesForAccount(ctx context.Context, id int64, cb GetActivitiesCallbackFunc) error {
	return nil
}

func (db *NullActivitiesDatabase) Close(ctx context.Context) error {
	return nil
}
