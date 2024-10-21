package queue

import (
	"context"
	"log/slog"
)

type SlogDeliveryQueue struct {
	DeliveryQueue
}

func init() {
	ctx := context.Background()
	RegisterDeliveryQueue(ctx, "slog", NewSlogDeliveryQueue)
}

func NewSlogDeliveryQueue(ctx context.Context, uri string) (DeliveryQueue, error) {
	q := &SlogDeliveryQueue{}
	return q, nil
}

func (q *SlogDeliveryQueue) DeliverActivity(ctx context.Context, opts *DeliverActivityOptions) error {
	slog.Info("Deliver post", "activity id", opts.Activity.Id, "from", opts.Activity.AccountId, "to", opts.To)
	return nil
}
