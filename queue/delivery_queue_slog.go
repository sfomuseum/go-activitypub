package queue

import (
	"context"
	"log/slog"

	"github.com/sfomuseum/go-activitypub"
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
	slog.Info("Deliver post", "post id", opts.Post.Id, "from", opts.From.Id, "to", opts.To)
	return nil
}
