package queue

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/aaronland/go-roster"
	"github.com/sfomuseum/go-activitypub/deliver"
)

type DeliveryQueue interface {
	DeliverActivity(context.Context, *deliver.DeliverActivityOptions) error
	Close(context.Context) error
}

var delivery_queue_roster roster.Roster

// DeliveryQueueInitializationFunc is a function defined by individual delivery_queue package and used to create
// an instance of that delivery_queue
type DeliveryQueueInitializationFunc func(ctx context.Context, uri string) (DeliveryQueue, error)

// RegisterDeliveryQueue registers 'scheme' as a key pointing to 'init_func' in an internal lookup table
// used to create new `DeliveryQueue` instances by the `NewDeliveryQueue` method.
func RegisterDeliveryQueue(ctx context.Context, scheme string, init_func DeliveryQueueInitializationFunc) error {

	err := ensureDeliveryQueueRoster()

	if err != nil {
		return err
	}

	return delivery_queue_roster.Register(ctx, scheme, init_func)
}

func ensureDeliveryQueueRoster() error {

	if delivery_queue_roster == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		delivery_queue_roster = r
	}

	return nil
}

// NewDeliveryQueue returns a new `DeliveryQueue` instance configured by 'uri'. The value of 'uri' is parsed
// as a `url.URL` and its scheme is used as the key for a corresponding `DeliveryQueueInitializationFunc`
// function used to instantiate the new `DeliveryQueue`. It is assumed that the scheme (and initialization
// function) have been registered by the `RegisterDeliveryQueue` method.
func NewDeliveryQueue(ctx context.Context, uri string) (DeliveryQueue, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := delivery_queue_roster.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	init_func := i.(DeliveryQueueInitializationFunc)
	return init_func(ctx, uri)
}

// Schemes returns the list of schemes that have been registered.
func DeliveryQueueSchemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureDeliveryQueueRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range delivery_queue_roster.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}
