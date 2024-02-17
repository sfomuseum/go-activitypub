package activitypub

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
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

func (q *SynchronousDeliveryQueue) DeliverPost(ctx context.Context, p *Post, follower_id string) error {

	to := []string{
		follower_id,
	}

	create_activity, err := p.AsCreateActivity(ctx, to)

	enc_activity, err := json.Marshal(create_activity)

	if err != nil {
		return fmt.Errorf("Failed to encode activity, %w", err)
	}

	slog.Info("POST", "activity", string(enc_activity))

	// Get followers inbox

	// Post to index

	// Happy happy

	return nil
}
