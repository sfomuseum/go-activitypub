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

type GetBoostIdsCallbackFunc func(context.Context, int64) error
type GetBoostsCallbackFunc func(context.Context, *activitypub.Boost) error

type BoostsDatabase interface {
	GetBoostIdsForDateRange(context.Context, int64, int64, GetBoostIdsCallbackFunc) error
	GetBoostsForPost(context.Context, int64, GetBoostsCallbackFunc) error
	GetBoostsForAccount(context.Context, int64, GetBoostsCallbackFunc) error
	// GetBoostsForActor(context.Context, string, GetBoostsCallbackFunc) error
	GetBoostWithPostIdAndActor(context.Context, int64, string) (*activitypub.Boost, error)
	GetBoostWithId(context.Context, int64) (*activitypub.Boost, error)
	AddBoost(context.Context, *activitypub.Boost) error
	RemoveBoost(context.Context, *activitypub.Boost) error
	Close(context.Context) error
}

var boost_database_roster roster.Roster

// BoostsDatabaseInitializationFunc is a function defined by individual boost_database package and used to create
// an instance of that boost_database
type BoostsDatabaseInitializationFunc func(ctx context.Context, uri string) (BoostsDatabase, error)

// RegisterBoostsDatabase registers 'scheme' as a key pointing to 'init_func' in an internal lookup table
// used to create new `BoostsDatabase` instances by the `NewBoostsDatabase` method.
func RegisterBoostsDatabase(ctx context.Context, scheme string, init_func BoostsDatabaseInitializationFunc) error {

	err := ensureBoostsDatabaseRoster()

	if err != nil {
		return err
	}

	return boost_database_roster.Register(ctx, scheme, init_func)
}

func ensureBoostsDatabaseRoster() error {

	if boost_database_roster == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		boost_database_roster = r
	}

	return nil
}

// NewBoostsDatabase returns a new `BoostsDatabase` instance configured by 'uri'. The value of 'uri' is parsed
// as a `url.URL` and its scheme is used as the key for a corresponding `BoostsDatabaseInitializationFunc`
// function used to instantiate the new `BoostsDatabase`. It is assumed that the scheme (and initialization
// function) have been registered by the `RegisterBoostsDatabase` method.
func NewBoostsDatabase(ctx context.Context, uri string) (BoostsDatabase, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := boost_database_roster.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	init_func := i.(BoostsDatabaseInitializationFunc)
	return init_func(ctx, uri)
}

// Schemes returns the list of schemes that have been registered.
func BoostsDatabaseSchemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureBoostsDatabaseRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range boost_database_roster.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}
