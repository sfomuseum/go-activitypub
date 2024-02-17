package activitypub

import (
	"context"
	"encoding/json"
	"fmt"

	"gocloud.dev/pubsub"
)

type PubSubDeliveryQueue struct {
	DeliveryQueue
	topic *pubsub.Topic
}

func init() {
	// ctx := context.Background()
	// RegisterDeliveryQueue(ctx, "pubSub", NewPubsubDeliveryQueue)
}

func NewPubSubDeliveryQueue(ctx context.Context, uri string) (DeliveryQueue, error) {

	topic, err := pubsub.OpenTopic(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to open pubsub topic, %w", err)
	}

	q := &PubSubDeliveryQueue{
		topic: topic,
	}
	return q, nil
}

func (q *PubSubDeliveryQueue) DeliverPost(ctx context.Context, opts *DeliverPostOptions) error {

	enc_opts, err := json.Marshal(opts)

	if err != nil {
		return fmt.Errorf("Failed to marshal post options, %w", err)
	}

	msg := &pubsub.Message{
		Body: enc_opts,
	}

	err = q.topic.Send(ctx, msg)

	if err != nil {
		return fmt.Errorf("Failed to send message, %w", err)
	}

	return nil
}
