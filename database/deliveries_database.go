package database

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/aaronland/go-roster"
	"github.com/sfomuseum/go-activitypub"
)

type GetDeliveryIdsCallbackFunc func(context.Context, int64) error
type GetDeliveriesCallbackFunc func(context.Context, *activitypub.Delivery) error

type DeliveriesDatabase interface {
	AddDelivery(context.Context, *activitypub.Delivery) error
	GetDeliveryWithId(context.Context, int64) (*activitypub.Delivery, error)
	GetDeliveries(context.Context, GetDeliveriesCallbackFunc) error
	GetDeliveriesWithActivityIdAndRecipient(context.Context, int64, string, GetDeliveriesCallbackFunc) error
	GetDeliveriesWithActivityPubIdAndRecipient(context.Context, string, string, GetDeliveriesCallbackFunc) error
	GetDeliveryIdsForDateRange(context.Context, int64, int64, GetDeliveryIdsCallbackFunc) error
	Close(context.Context) error
}

var deliveries_database_roster roster.Roster

// DeliveriesDatabaseInitializationFunc is a function defined by individual deliveries_database package and used to create
// an instance of that deliveries_database
type DeliveriesDatabaseInitializationFunc func(ctx context.Context, uri string) (DeliveriesDatabase, error)

// RegisterDeliveriesDatabase registers 'scheme' as a key pointing to 'init_func' in an internal lookup table
// used to create new `DeliveriesDatabase` instances by the `NewDeliveriesDatabase` method.
func RegisterDeliveriesDatabase(ctx context.Context, scheme string, init_func DeliveriesDatabaseInitializationFunc) error {

	err := ensureDeliveriesDatabaseRoster()

	if err != nil {
		return err
	}

	return deliveries_database_roster.Register(ctx, scheme, init_func)
}

func ensureDeliveriesDatabaseRoster() error {

	if deliveries_database_roster == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		deliveries_database_roster = r
	}

	return nil
}

// NewDeliveriesDatabase returns a new `DeliveriesDatabase` instance configured by 'uri'. The value of 'uri' is parsed
// as a `url.URL` and its scheme is used as the key for a corresponding `DeliveriesDatabaseInitializationFunc`
// function used to instantiate the new `DeliveriesDatabase`. It is assumed that the scheme (and initialization
// function) have been registered by the `RegisterDeliveriesDatabase` method.
func NewDeliveriesDatabase(ctx context.Context, uri string) (DeliveriesDatabase, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := deliveries_database_roster.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	init_func := i.(DeliveriesDatabaseInitializationFunc)
	return init_func(ctx, uri)
}

// Schemes returns the list of schemes that have been registered.
func DeliveriesDatabaseSchemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureDeliveriesDatabaseRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range deliveries_database_roster.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}
