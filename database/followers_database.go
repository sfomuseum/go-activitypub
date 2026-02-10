package database

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"sort"
	"strings"
	"sync/atomic"
	"time"

	"github.com/aaronland/go-roster"
	"github.com/sfomuseum/go-activitypub"
)

type GetFollowerIdsCallbackFunc func(context.Context, int64) error
type GetFollowersCallbackFunc func(context.Context, string) error
type GetFollowersCallbackFunc2 func(context.Context, *activitypub.Follower) error

type FollowersDatabase interface {
	GetFollower(context.Context, int64, string) (*activitypub.Follower, error)
	AddFollower(context.Context, *activitypub.Follower) error
	RemoveFollower(context.Context, *activitypub.Follower) error
	GetFollowers(context.Context, GetFollowersCallbackFunc2) error // Get all the follower rows
	GetFollowerWithId(context.Context, int64) (*activitypub.Follower, error)
	GetFollowerIdsForDateRange(context.Context, int64, int64, GetFollowerIdsCallbackFunc) error
	GetAllFollowers(context.Context, GetFollowersCallbackFunc) error // Get all follower addresses (probably deprecated)
	GetFollowersForAccount(context.Context, int64, GetFollowersCallbackFunc) error
	HasFollowers(context.Context, int64) (bool, error)
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

func MigrateFollowersDatabaseFromURIs(ctx context.Context, from_uri string, to_uri string, count *int64, success *int64, errors *int64) error {

	from_ctx, from_cancel := context.WithTimeout(ctx, 5*time.Second)
	defer from_cancel()

	from_db, err := NewFollowersDatabase(from_ctx, from_uri)

	if err != nil {
		return fmt.Errorf("Failed to create from database, %w", err)
	}

	defer from_db.Close(ctx)

	slog.Debug("Set up to database")

	to_ctx, to_cancel := context.WithTimeout(ctx, 5*time.Second)
	defer to_cancel()

	to_db, err := NewFollowersDatabase(to_ctx, to_uri)

	if err != nil {
		return fmt.Errorf("Failed to create to database, %w", err)
	}

	defer to_db.Close(ctx)

	return MigrateFollowersDatabase(ctx, from_db, to_db, count, success, errors)
}

func MigrateFollowersDatabase(ctx context.Context, from_db FollowersDatabase, to_db FollowersDatabase, count *int64, success *int64, errors *int64) error {

	cb := func(ctx context.Context, d *activitypub.Follower) error {

		defer atomic.AddInt64(count, 1)

		slog.Debug("Add", "follower", d.Id)
		err := to_db.AddFollower(ctx, d)

		if err != nil {
			slog.Error("Failed to add follower", "follower", d.Id, "error", err)
			atomic.AddInt64(errors, 1)
		} else {
			atomic.AddInt64(success, 1)
		}

		return nil
	}

	slog.Debug("Retrieve followers")
	return from_db.GetFollowers(ctx, cb)
}
