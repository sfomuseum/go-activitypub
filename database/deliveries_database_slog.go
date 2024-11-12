package database

import (
	"context"
	"log/slog"

	"github.com/sfomuseum/go-activitypub"
)

type SlogDeliveriesDatabase struct {
	DeliveriesDatabase
	logger *slog.Logger
}

func init() {
	ctx := context.Background()
	err := RegisterDeliveriesDatabase(ctx, "slog", NewSlogDeliveriesDatabase)

	if err != nil {
		panic(err)
	}
}

func NewSlogDeliveriesDatabase(ctx context.Context, uri string) (DeliveriesDatabase, error) {
	db := &SlogDeliveriesDatabase{
		logger: slog.Default(),
	}
	return db, nil
}

func (db *SlogDeliveriesDatabase) AddDelivery(ctx context.Context, d *activitypub.Delivery) error {
	db.logger.Info("Add delivery", "activity id", d.ActivityId, "recipient", d.Recipient, "success", d.Success, "error", d.Error)
	return nil
}

func (db *SlogDeliveriesDatabase) GetDeliveryWithId(ctx context.Context, id int64) (*activitypub.Delivery, error) {
	return nil, activitypub.ErrNotFound
}

func (db *SlogDeliveriesDatabase) GetDeliveries(ctx context.Context, cb GetDeliveriesCallbackFunc) error {
	return nil
}

func (db *SlogDeliveriesDatabase) GetDeliveriesWithActivityIdAndRecipient(ctx context.Context, activity_id int64, recipient string, cb GetDeliveriesCallbackFunc) error {
	return nil
}

func (db *SlogDeliveriesDatabase) GetDeliveriesWithActivityPubIdAndRecipient(ctx context.Context, activity_pub_id string, recipient string, cb GetDeliveriesCallbackFunc) error {
	return nil
}

func (db *SlogDeliveriesDatabase) GetDeliveryIdsForDateRange(ctx context.Context, start int64, end int64, cb GetDeliveryIdsCallbackFunc) error {
	return nil
}

func (db *SlogDeliveriesDatabase) Close(ctx context.Context) error {
	return nil
}
