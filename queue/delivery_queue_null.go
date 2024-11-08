package queue

import (
	"context"

	"github.com/sfomuseum/go-activitypub/deliver"
)

type NullDeliveryQueue struct {
	DeliveryQueue
}

func init() {
	ctx := context.Background()
	err := RegisterDeliveryQueue(ctx, "null", NewNullDeliveryQueue)

	if err != nil {
		panic(err)
	}
}

func NewNullDeliveryQueue(ctx context.Context, uri string) (DeliveryQueue, error) {
	q := &NullDeliveryQueue{}
	return q, nil
}

func (q *NullDeliveryQueue) DeliverActivity(ctx context.Context, opts *deliver.DeliverActivityOptions) error {
	return nil
}

func (q *NullDeliveryQueue) Close(ctx context.Context) error {
	return nil
}
