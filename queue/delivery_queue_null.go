package activitypub

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

func (q *NullDeliveryQueue) DeliverPost(ctx context.Context, opts *DeliverPostOptions) error {
	return nil
}
