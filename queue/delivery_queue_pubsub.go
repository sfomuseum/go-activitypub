package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/sfomuseum/go-activitypub/deliver"
	"github.com/sfomuseum/go-pubsub/publisher"
)

type PubSubDeliveryQueueOptions struct {
	// The actor to whom the activity should be delivered.
	To string `json:"to"`
	// Remember PostId is a misnomer. See notes in activity.go
	ActivityId int64 `json:"activity_id"`
}

type PubSubDeliveryQueue struct {
	DeliveryQueue
	publisher publisher.Publisher
}

// In principle this could also be done with a sync.OnceFunc call but that will
// require that everyone uses Go 1.21 (whose package import changes broke everything)
// which is literally days old as I write this. So maybe a few releases after 1.21.
//
// Also, _not_ using a sync.OnceFunc means we can call RegisterSchemes multiple times
// if and when multiple gomail-sender instances register themselves.

var register_mu = new(sync.RWMutex)
var register_map = map[string]bool{}

func init() {

	ctx := context.Background()
	err := RegisterPubSubSchemes(ctx)

	if err != nil {
		panic(err)
	}
}

// RegisterSchemes will explicitly register all the schemes associated with the `AccessTokensDeliveryAgent` interface.
func RegisterPubSubSchemes(ctx context.Context) error {

	register_mu.Lock()
	defer register_mu.Unlock()

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

		_, exists := register_map[scheme]

		if exists {
			continue
		}

		err := RegisterDeliveryQueue(ctx, scheme, NewPubSubDeliveryQueue)

		if err != nil {
			return fmt.Errorf("Failed to register delivery queue for '%s', %w", scheme, err)
		}

		register_map[scheme] = true
	}

	return nil
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

func (q *PubSubDeliveryQueue) DeliverActivity(ctx context.Context, opts *deliver.DeliverActivityOptions) error {

	ps_opts := PubSubDeliveryQueueOptions{
		To:         opts.To,
		ActivityId: opts.Activity.Id,
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
