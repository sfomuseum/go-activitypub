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
	RegisterDeliveryQueue(ctx, "synchronous", NewSynchronousDeliveryQueue)
}

func NewSynchronousDeliveryQueue(ctx context.Context, uri string) (DeliveryQueue, error) {
	q := &SynchronousDeliveryQueue{}
	return q, nil
}

func (q *SynchronousDeliveryQueue) DeliverActivity(ctx context.Context, opts *DeliverActivityOptions) error {

	err := deliver.DeliverActivity(ctx, opts)

	if err != nil {
		return fmt.Errorf("Failed to deliver post, %w", err)
	}

	return nil
}
