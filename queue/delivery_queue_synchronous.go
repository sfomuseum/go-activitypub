package queue

import (
	"context"
	"fmt"

	"github.com/sfomuseum/go-activitypub/deliver"
)

type SynchronousDeliveryQueue struct {
	DeliveryQueue
}

func init() {
	ctx := context.Background()
	err := RegisterDeliveryQueue(ctx, "synchronous", NewSynchronousDeliveryQueue)

	if err != nil {
		panic(err)
	}
}

func NewSynchronousDeliveryQueue(ctx context.Context, uri string) (DeliveryQueue, error) {
	q := &SynchronousDeliveryQueue{}
	return q, nil
}

func (q *SynchronousDeliveryQueue) DeliverActivity(ctx context.Context, opts *deliver.DeliverActivityOptions) error {

	err := deliver.DeliverActivity(ctx, opts)

	if err != nil {
		return fmt.Errorf("Failed to deliver post, %w", err)
	}

	return nil
}

func (q *SynchronousDeliveryQueue) Close(ctx context.Context) error {
	return nil
}
