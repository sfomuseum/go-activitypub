package activitypub

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sfomuseum/go-pubsub/publisher"
)

type PubSubDeliveryQueuePostOptions struct {
	AccountId int64  `json:"account_id"`
	Recipient string `json:"recipient"`
	PostId    int64  `json:"post_id"`
}

type PubSubDeliveryQueue struct {
	DeliveryQueue
	publisher publisher.Publisher
}

func init() {

	ctx := context.Background()

	to_register := []string{
		"awssqs-creds",
	}

	for _, scheme := range to_register {
		RegisterDeliveryQueue(ctx, scheme, NewPubSubDeliveryQueue)
	}

	for _, scheme := range publisher.PublisherSchemes() {
		RegisterDeliveryQueue(ctx, scheme, NewPubSubDeliveryQueue)
	}
}

func NewPubSubDeliveryQueue(ctx context.Context, uri string) (DeliveryQueue, error) {

	pub, err := publisher.NewPublisher(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create publisher, %w", err)
	}

	q := &PubSubDeliveryQueue{
		publisher: pub,
	}

	return q, nil
}

func (q *PubSubDeliveryQueue) DeliverPost(ctx context.Context, opts *DeliverPostOptions) error {

	ps_opts := PubSubDeliveryQueuePostOptions{
		AccountId: opts.From.Id,
		Recipient: opts.To,
		PostId:    opts.Post.Id,
	}

	enc_opts, err := json.Marshal(ps_opts)

	if err != nil {
		return fmt.Errorf("Failed to marshal post options, %w", err)
	}

	err = q.publisher.Publish(ctx, string(enc_opts))

	if err != nil {
		return fmt.Errorf("Failed to send message, %w", err)
	}

	return nil
}
