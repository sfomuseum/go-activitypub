package queue

import (
	"context"
	"log/slog"

	"github.com/sfomuseum/go-activitypub/deliver"
)

type SlogDeliveryQueue struct {
	DeliveryQueue
}

func init() {
	ctx := context.Background()
	err := RegisterDeliveryQueue(ctx, "slog", NewSlogDeliveryQueue)

	if err != nil {
		panic(err)
	}
}

func NewSlogDeliveryQueue(ctx context.Context, uri string) (DeliveryQueue, error) {
	q := &SlogDeliveryQueue{}
	return q, nil
}

func (q *SlogDeliveryQueue) DeliverActivity(ctx context.Context, opts *deliver.DeliverActivityOptions) error {
	slog.Info("Deliver post", "activity id", opts.Activity.Id, "from", opts.Activity.AccountId, "to", opts.To)
	return nil
}

func (q *SlogDeliveryQueue) Close(ctx context.Context) error {
	return nil
}
