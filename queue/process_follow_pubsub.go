package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/sfomuseum/go-pubsub/publisher"
)

type PubSubProcessFollowQueue struct {
	ProcessFollowQueue
	publisher publisher.Publisher
}

var process_follow_mu = new(sync.RWMutex)
var process_follow_map = map[string]bool{}

func init() {

	ctx := context.Background()

	err := RegisterPubSubProcessFollowSchemes(ctx)

	if err != nil {
		panic(err)
	}
}

func RegisterPubSubProcessFollowSchemes(ctx context.Context) error {

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

		err := RegisterProcessFollowQueue(ctx, scheme, NewPubSubProcessFollowQueue)

		if err != nil {
			return fmt.Errorf("Failed to register delivery queue for '%s', %w", scheme, err)
		}

		process_follow_map[scheme] = true
	}

	return nil
}

func NewPubSubProcessFollowQueue(ctx context.Context, uri string) (ProcessFollowQueue, error) {

	pub, err := publisher.NewPublisher(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create publisher, %w", err)
	}

	q := &PubSubProcessFollowQueue{
		publisher: pub,
	}

	return q, nil
}

func (q *PubSubProcessFollowQueue) ProcessFollow(ctx context.Context, follower_id int64) error {

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

func (q *PubSubProcessFollowQueue) Close(ctx context.Context) error {
	return q.publisher.Close()
}
