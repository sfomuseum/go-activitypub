package activitypub

import (
	"context"
	"fmt"
)

type SynchronousDeliveryQueue struct {
	DeliveryQueue
}

func init() {
	ctx := context.Background()
	RegisterDeliveryQueue(ctx, "synchronous", NewSynchronousDeliveryQueue)
}

func NewSynchronousDeliveryQueue(ctx context.Context, uri string) (DeliveryQueue, error) {
	q := &SynchronousDeliveryQueue{}
	return q, nil
}

func (q *SynchronousDeliveryQueue) DeliverPost(ctx context.Context, opts *DeliverPostOptions) error {

	err := DeliverPost(ctx, opts)

	if err != nil {
		return fmt.Errorf("Failed to deliver post, %w", err)
	}

	return nil
}
