package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/sfomuseum/go-pubsub/publisher"
)

type PubSubProcessMessageQueue struct {
	ProcessMessageQueue
	publisher publisher.Publisher
}

var process_register_mu = new(sync.RWMutex)
var process_register_map = map[string]bool{}

func init() {

	ctx := context.Background()

	err := RegisterPubSubProcessMessageSchemes(ctx)

	if err != nil {
		panic(err)
	}
}

func RegisterPubSubProcessMessageSchemes(ctx context.Context) error {

	process_register_mu.Lock()
	defer process_register_mu.Unlock()

	to_register := []string{
		"awssqs-creds",
	}

	for _, scheme := range publisher.PublisherSchemes() {

		scheme = strings.Replace(scheme, "://", "", 1)

		// I don't love this so maybe prefix everything as pubsub-{SCHEME} ? TBD... ?

		if scheme != "null" {
			to_register = append(to_register, scheme)
		}
	}

	for _, scheme := range to_register {

		_, exists := process_register_map[scheme]

		if exists {
			continue
		}

		err := RegisterProcessMessageQueue(ctx, scheme, NewPubSubProcessMessageQueue)

		if err != nil {
			return fmt.Errorf("Failed to register delivery queue for '%s', %w", scheme, err)
		}

		process_register_map[scheme] = true
	}

	return nil
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

func (q *PubSubProcessMessageQueue) Close(ctx context.Context) error {
	return q.publisher.Close()
}
