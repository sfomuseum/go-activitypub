package activitypub

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/aaronland/go-roster"
)

type ActorDatabase interface {
	GetActor(context.Context, string) (*Actor, error)
	AddActor(context.Context, *Actor) error
	RemoveActor(context.Context, *Actor) error
	UpdateActor(context.Context, *Actor) error
}

var actor_database_roster roster.Roster

// ActorDatabaseInitializationFunc is a function defined by individual actor_database package and used to create
// an instance of that actor_database
type ActorDatabaseInitializationFunc func(ctx context.Context, uri string) (ActorDatabase, error)

// RegisterActorDatabase registers 'scheme' as a key pointing to 'init_func' in an internal lookup table
// used to create new `ActorDatabase` instances by the `NewActorDatabase` method.
func RegisterActorDatabase(ctx context.Context, scheme string, init_func ActorDatabaseInitializationFunc) error {

	err := ensureActorDatabaseRoster()

	if err != nil {
		return err
	}

	return actor_database_roster.Register(ctx, scheme, init_func)
}

func ensureActorDatabaseRoster() error {

	if actor_database_roster == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		actor_database_roster = r
	}

	return nil
}

// NewActorDatabase returns a new `ActorDatabase` instance configured by 'uri'. The value of 'uri' is parsed
// as a `url.URL` and its scheme is used as the key for a corresponding `ActorDatabaseInitializationFunc`
// function used to instantiate the new `ActorDatabase`. It is assumed that the scheme (and initialization
// function) have been registered by the `RegisterActorDatabase` method.
func NewActorDatabase(ctx context.Context, uri string) (ActorDatabase, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := actor_database_roster.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	init_func := i.(ActorDatabaseInitializationFunc)
	return init_func(ctx, uri)
}

// Schemes returns the list of schemes that have been registered.
func ActorDatabaseSchemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureActorDatabaseRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range actor_database_roster.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}
