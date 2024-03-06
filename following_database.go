package activitypub

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/aaronland/go-roster"
)

type GetFollowingIdsCallbackFunc func(context.Context, int64) error
type GetFollowingCallbackFunc func(context.Context, string) error

type FollowingDatabase interface {
	GetFollowingIdsForDateRange(context.Context, int64, int64, GetFollowingIdsCallbackFunc) error
	GetFollowingForAccount(context.Context, int64, GetFollowingCallbackFunc) error
	GetFollowing(context.Context, int64, string) (*Following, error)
	AddFollowing(context.Context, *Following) error
	RemoveFollowing(context.Context, *Following) error
	Close(context.Context) error
}

var following_database_roster roster.Roster

// FollowingDatabaseInitializationFunc is a function defined by individual following_database package and used to create
// an instance of that following_database
type FollowingDatabaseInitializationFunc func(ctx context.Context, uri string) (FollowingDatabase, error)

// RegisterFollowingDatabase registers 'scheme' as a key pointing to 'init_func' in an internal lookup table
// used to create new `FollowingDatabase` instances by the `NewFollowingDatabase` method.
func RegisterFollowingDatabase(ctx context.Context, scheme string, init_func FollowingDatabaseInitializationFunc) error {

	err := ensureFollowingDatabaseRoster()

	if err != nil {
		return err
	}

	return following_database_roster.Register(ctx, scheme, init_func)
}

func ensureFollowingDatabaseRoster() error {

	if following_database_roster == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		following_database_roster = r
	}

	return nil
}

// NewFollowingDatabase returns a new `FollowingDatabase` instance configured by 'uri'. The value of 'uri' is parsed
// as a `url.URL` and its scheme is used as the key for a corresponding `FollowingDatabaseInitializationFunc`
// function used to instantiate the new `FollowingDatabase`. It is assumed that the scheme (and initialization
// function) have been registered by the `RegisterFollowingDatabase` method.
func NewFollowingDatabase(ctx context.Context, uri string) (FollowingDatabase, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := following_database_roster.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	init_func := i.(FollowingDatabaseInitializationFunc)
	return init_func(ctx, uri)
}

// Schemes returns the list of schemes that have been registered.
func FollowingDatabaseSchemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureFollowingDatabaseRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range following_database_roster.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}
