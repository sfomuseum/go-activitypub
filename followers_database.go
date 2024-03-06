package activitypub

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/aaronland/go-roster"
)

type GetFollowerIdsCallbackFunc func(context.Context, int64) error
type GetFollowersCallbackFunc func(context.Context, string) error

type FollowersDatabase interface {
	GetFollowerIdsForDateRange(context.Context, int64, int64, GetFollowerIdsCallbackFunc) error
	GetFollowersForAccount(context.Context, int64, GetFollowersCallbackFunc) error
	HasFollowers(context.Context, int64) (bool, error)
	GetFollower(context.Context, int64, string) (*Follower, error)
	AddFollower(context.Context, *Follower) error
	RemoveFollower(context.Context, *Follower) error
	Close(context.Context) error
}

var followers_database_roster roster.Roster

// FollowersDatabaseInitializationFunc is a function defined by individual followers_database package and used to create
// an instance of that followers_database
type FollowersDatabaseInitializationFunc func(ctx context.Context, uri string) (FollowersDatabase, error)

// RegisterFollowersDatabase registers 'scheme' as a key pointing to 'init_func' in an internal lookup table
// used to create new `FollowersDatabase` instances by the `NewFollowersDatabase` method.
func RegisterFollowersDatabase(ctx context.Context, scheme string, init_func FollowersDatabaseInitializationFunc) error {

	err := ensureFollowersDatabaseRoster()

	if err != nil {
		return err
	}

	return followers_database_roster.Register(ctx, scheme, init_func)
}

func ensureFollowersDatabaseRoster() error {

	if followers_database_roster == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		followers_database_roster = r
	}

	return nil
}

// NewFollowersDatabase returns a new `FollowersDatabase` instance configured by 'uri'. The value of 'uri' is parsed
// as a `url.URL` and its scheme is used as the key for a corresponding `FollowersDatabaseInitializationFunc`
// function used to instantiate the new `FollowersDatabase`. It is assumed that the scheme (and initialization
// function) have been registered by the `RegisterFollowersDatabase` method.
func NewFollowersDatabase(ctx context.Context, uri string) (FollowersDatabase, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := followers_database_roster.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	init_func := i.(FollowersDatabaseInitializationFunc)
	return init_func(ctx, uri)
}

// Schemes returns the list of schemes that have been registered.
func FollowersDatabaseSchemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureFollowersDatabaseRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range followers_database_roster.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}
