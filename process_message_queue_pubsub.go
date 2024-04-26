package activitypub

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sfomuseum/go-pubsub/publisher"
)

type PubSubProcessMessageQueue struct {
	ProcessMessageQueue
	publisher publisher.Publisher
}

func init() {

	ctx := context.Background()

	to_register := []string{
		"awssqs-creds",
	}

	for _, scheme := range to_register {
		RegisterProcessMessageQueue(ctx, scheme, NewPubSubProcessMessageQueue)
	}

	for _, scheme := range publisher.PublisherSchemes() {
		scheme = strings.Replace(scheme, "://", "", 1)
		RegisterProcessMessageQueue(ctx, scheme, NewPubSubProcessMessageQueue)
	}
}

func NewPubSubProcessMessageQueue(ctx context.Context, uri string) (ProcessMessageQueue, error) {

	pub, err := publisher.NewPublisher(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create publisher, %w", err)
	}

	q := &PubSubProcessMessageQueue{
		publisher: pub,
	}

	return q, nil
}

func (q *PubSubProcessMessageQueue) ProcessMessage(ctx context.Context, message_id int64) error {

	enc_id, err := json.Marshal(message_id)

	if err != nil {
		return fmt.Errorf("Failed to marshal message ID, %w", err)
	}

	err = q.publisher.Publish(ctx, string(enc_id))

	if err != nil {
		return fmt.Errorf("Failed to send message, %w", err)
	}

	return nil
}
