package activitypub

import (
	"context"
)

type NullPropertiesDatabase struct {
	PropertiesDatabase
}

func init() {
	ctx := context.Background()
	RegisterPropertiesDatabase(ctx, "null", NewNullPropertiesDatabase)
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

func (db *NullPropertiesDatabase) AddProperty(ctx context.Context, property *Property) error {
	return nil
}

func (db *NullPropertiesDatabase) UpdateProperty(ctx context.Context, property *Property) error {
	return nil
}

func (db *NullPropertiesDatabase) RemoveProperty(ctx context.Context, property *Property) error {
	return nil
}

func (db *NullPropertiesDatabase) Close(ctx context.Context) error {
	return nil
}
