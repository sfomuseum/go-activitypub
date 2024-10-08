package activitypub

import (
	"context"
	"log/slog"
)

type SlogDeliveriesDatabase struct {
	DeliveriesDatabase
}

func init() {
	ctx := context.Background()
	RegisterDeliveriesDatabase(ctx, "slog", NewSlogDeliveriesDatabase)
}

func NewSlogDeliveriesDatabase(ctx context.Context, uri string) (DeliveriesDatabase, error) {
	db := &SlogDeliveriesDatabase{}
	return db, nil
}

func (db *SlogDeliveriesDatabase) AddDelivery(ctx context.Context, d *Delivery) error {
	slog.Info("Add delivery", "post id", d.PostId, "recipient", d.Recipient, "success", d.Success, "error", d.Error)
	return nil
}

func (db *SlogDeliveriesDatabase) GetDeliveryWithId(ctx context.Context, id int64) (*Delivery, error) {
	return nil, ErrNotFound
}

func (db *SlogDeliveriesDatabase) GetDeliveriesWithPostIdAndRecipient(ctx context.Context, post_id int64, recipient string, cb GetDeliveriesCallbackFunc) error {
	return nil
}

func (db *SlogDeliveriesDatabase) Close(ctx context.Context) error {
	return nil
}
