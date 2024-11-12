package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"github.com/sfomuseum/go-activitypub/deliver"
	"github.com/sfomuseum/go-activitypub/id"
	"github.com/sfomuseum/go-pubsub/publisher"
)

type PubSubDeliveryQueueOptions struct {
	// The unique ID associated with the pubsub delivery. This is mostly for debugging between the sender and the receiver.
	Id int64 `json:"id"`
	// The actor to whom the activity should be delivered.
	To string `json:"to"`
	// The unique Activity(Database) Id associated with the delivery.
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

var delivery_register_mu = new(sync.RWMutex)
var delivery_register_map = map[string]bool{}

func init() {

	ctx := context.Background()
	err := RegisterPubSubDeliverySchemes(ctx)

	if err != nil {
		panic(err)
	}
}

func RegisterPubSubDeliverySchemes(ctx context.Context) error {

	delivery_register_mu.Lock()

	defer delivery_register_mu.Unlock()

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

		_, exists := delivery_register_map[scheme]

		if exists {
			continue
		}

		err := RegisterDeliveryQueue(ctx, scheme, NewPubSubDeliveryQueue)

		if err != nil {
			return fmt.Errorf("Failed to register delivery queue for '%s', %w", scheme, err)
		}

		delivery_register_map[scheme] = true
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

	ps_id, err := id.NewId()

	if err != nil {
		return fmt.Errorf("Failed to create unique pubsub ID, %w", err)
	}

	ps_opts := PubSubDeliveryQueueOptions{
		Id:         ps_id,
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

	logger := slog.Default()
	logger.Info("Published pubsub activity", "to", opts.To, "from", opts.Activity.AccountId, "activity id", opts.Activity.Id, "pubsub id", ps_id)

	return nil
}

func (q *PubSubDeliveryQueue) Close(ctx context.Context) error {
	return q.publisher.Close()
}
