package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/sfomuseum/go-pubsub/publisher"
)

type PubSubProcessFollowerQueue struct {
	ProcessFollowerQueue
	publisher publisher.Publisher
}

var process_follow_mu = new(sync.RWMutex)
var process_follow_map = map[string]bool{}

func init() {

	ctx := context.Background()

	err := RegisterPubSubProcessFollowerSchemes(ctx)

	if err != nil {
		panic(err)
	}
}

func RegisterPubSubProcessFollowerSchemes(ctx context.Context) error {

	process_follow_mu.Lock()
	defer process_follow_mu.Unlock()

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

		_, exists := process_follow_map[scheme]

		if exists {
			continue
		}

		err := RegisterProcessFollowerQueue(ctx, scheme, NewPubSubProcessFollowerQueue)

		if err != nil {
			return fmt.Errorf("Failed to register delivery queue for '%s', %w", scheme, err)
		}

		process_follow_map[scheme] = true
	}

	return nil
}

func NewPubSubProcessFollowerQueue(ctx context.Context, uri string) (ProcessFollowerQueue, error) {

	pub, err := publisher.NewPublisher(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create publisher, %w", err)
	}

	q := &PubSubProcessFollowerQueue{
		publisher: pub,
	}

	return q, nil
}

func (q *PubSubProcessFollowerQueue) ProcessFollower(ctx context.Context, follower_id int64) error {

	enc_msg, err := json.Marshal(follower_id)

	if err != nil {
		return fmt.Errorf("Failed to marshal message, %w", err)
	}

	err = q.publisher.Publish(ctx, string(enc_msg))

	if err != nil {
		return fmt.Errorf("Failed to send message, %w", err)
	}

	return nil
}

func (q *PubSubProcessFollowerQueue) Close(ctx context.Context) error {
	return q.publisher.Close()
}
