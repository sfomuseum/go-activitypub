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

type GetLikeIdsCallbackFunc func(context.Context, int64) error
type GetLikesCallbackFunc func(context.Context, *activitypub.Like) error

type LikesDatabase interface {
	GetLikeIdsForDateRange(context.Context, int64, int64, GetLikeIdsCallbackFunc) error
	GetLikesForPost(context.Context, int64, GetLikesCallbackFunc) error
	GetLikeWithPostIdAndActor(context.Context, int64, string) (*activitypub.Like, error)
	GetLikeWithId(context.Context, int64) (*activitypub.Like, error)
	AddLike(context.Context, *activitypub.Like) error
	RemoveLike(context.Context, *activitypub.Like) error
	Close(context.Context) error
}

var like_database_roster roster.Roster

// LikesDatabaseInitializationFunc is a function defined by individual like_database package and used to create
// an instance of that like_database
type LikesDatabaseInitializationFunc func(ctx context.Context, uri string) (LikesDatabase, error)

// RegisterLikesDatabase registers 'scheme' as a key pointing to 'init_func' in an internal lookup table
// used to create new `LikesDatabase` instances by the `NewLikesDatabase` method.
func RegisterLikesDatabase(ctx context.Context, scheme string, init_func LikesDatabaseInitializationFunc) error {

	err := ensureLikesDatabaseRoster()

	if err != nil {
		return err
	}

	return like_database_roster.Register(ctx, scheme, init_func)
}

func ensureLikesDatabaseRoster() error {

	if like_database_roster == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		like_database_roster = r
	}

	return nil
}

// NewLikesDatabase returns a new `LikesDatabase` instance configured by 'uri'. The value of 'uri' is parsed
// as a `url.URL` and its scheme is used as the key for a corresponding `LikesDatabaseInitializationFunc`
// function used to instantiate the new `LikesDatabase`. It is assumed that the scheme (and initialization
// function) have been registered by the `RegisterLikesDatabase` method.
func NewLikesDatabase(ctx context.Context, uri string) (LikesDatabase, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := like_database_roster.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	init_func := i.(LikesDatabaseInitializationFunc)
	return init_func(ctx, uri)
}

// Schemes returns the list of schemes that have been registered.
func LikesDatabaseSchemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureLikesDatabaseRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range like_database_roster.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}
