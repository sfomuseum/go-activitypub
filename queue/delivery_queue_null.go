package queue

import (
	"context"
)

type NullDeliveryQueue struct {
	DeliveryQueue
}

func init() {
	ctx := context.Background()
	RegisterDeliveryQueue(ctx, "null", NewNullDeliveryQueue)
}

func NewNullDeliveryQueue(ctx context.Context, uri string) (DeliveryQueue, error) {
	q := &NullDeliveryQueue{}
	return q, nil
}

func (q *NullDeliveryQueue) DeliverActivity(ctx context.Context, opts *DeliverActivityOptions) error {
	return nil
}