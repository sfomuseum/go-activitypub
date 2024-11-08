package database

import (
	"context"

	"github.com/sfomuseum/go-activitypub"
)

type NullPropertiesDatabase struct {
	PropertiesDatabase
}

func init() {
	ctx := context.Background()
	err := RegisterPropertiesDatabase(ctx, "null", NewNullPropertiesDatabase)

	if err != nil {
		panic(err)
	}
}

func NewNullPropertiesDatabase(ctx context.Context, uri string) (PropertiesDatabase, error) {
	db := &NullPropertiesDatabase{}
	return db, nil
}

func (db *NullPropertiesDatabase) GetProperties(ctx context.Context, cb GetPropertiesCallbackFunc) error {
	return nil
}

func (db *NullPropertiesDatabase) GetPropertiesForAccount(ctx context.Context, account_id int64, cb GetPropertiesCallbackFunc) error {
	return nil
}

func (db *NullPropertiesDatabase) AddProperty(ctx context.Context, property *activitypub.Property) error {
	return nil
}

func (db *NullPropertiesDatabase) UpdateProperty(ctx context.Context, property *activitypub.Property) error {
	return nil
}

func (db *NullPropertiesDatabase) RemoveProperty(ctx context.Context, property *activitypub.Property) error {
	return nil
}

func (db *NullPropertiesDatabase) Close(ctx context.Context) error {
	return nil
}
