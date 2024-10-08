package database

import (
	"context"

	"github.com/sfomuseum/go-activitypub"
)

type NullDeliveriesDatabase struct {
	DeliveriesDatabase
}

func init() {
	ctx := context.Background()
	RegisterDeliveriesDatabase(ctx, "null", NewNullDeliveriesDatabase)
}

func NewNullDeliveriesDatabase(ctx context.Context, uri string) (DeliveriesDatabase, error) {
	db := &NullDeliveriesDatabase{}
	return db, nil
}

func (db *NullDeliveriesDatabase) AddDelivery(ctx context.Context, d *activitypub.Delivery) error {
	return nil
}

func (db *NullDeliveriesDatabase) GetDeliveryWithId(ctx context.Context, id int64) (*activitypub.Delivery, error) {
	return nil, ErrNotFound
}

func (db *NullDeliveriesDatabase) GetDeliveriesWithPostIdAndRecipient(ctx context.Context, post_id int64, recipient string, cb GetDeliveriesCallbackFunc) error {
	return nil
}

func (db *NullDeliveriesDatabase) Close(ctx context.Context) error {
	return nil
}
