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
	err := RegisterDeliveriesDatabase(ctx, "null", NewNullDeliveriesDatabase)

	if err != nil {
		panic(err)
	}
}

func NewNullDeliveriesDatabase(ctx context.Context, uri string) (DeliveriesDatabase, error) {
	db := &NullDeliveriesDatabase{}
	return db, nil
}

func (db *NullDeliveriesDatabase) AddDelivery(ctx context.Context, d *activitypub.Delivery) error {
	return nil
}

func (db *NullDeliveriesDatabase) GetDeliveryWithId(ctx context.Context, id int64) (*activitypub.Delivery, error) {
	return nil, activitypub.ErrNotFound
}

func (db *NullDeliveriesDatabase) GetDeliveries(ctx context.Context, cb GetDeliveriesCallbackFunc) error {
	return nil
}

func (db *NullDeliveriesDatabase) GetDeliveriesWithActivityIdAndRecipient(ctx context.Context, activity_id int64, recipient string, cb GetDeliveriesCallbackFunc) error {
	return nil
}

func (db *NullDeliveriesDatabase) GetDeliveriesWithActivityPubIdAndRecipient(ctx context.Context, activity_pub_id string, recipient string, cb GetDeliveriesCallbackFunc) error {
	return nil
}

func (db *NullDeliveriesDatabase) GetDeliveryIdsForDateRange(ctx context.Context, start int64, end int64, cb GetDeliveryIdsCallbackFunc) error {
	return nil
}

func (db *NullDeliveriesDatabase) Close(ctx context.Context) error {
	return nil
}
