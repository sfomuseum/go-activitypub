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

type GetPostTagIdsCallbackFunc func(context.Context, int64) error
type GetPostTagsCallbackFunc func(context.Context, *activitypub.PostTag) error

type PostTagsDatabase interface {
	GetPostTagIdsForDateRange(context.Context, int64, int64, GetPostTagIdsCallbackFunc) error
	GetPostTagsForName(context.Context, string, GetPostTagsCallbackFunc) error
	GetPostTagsForAccount(context.Context, int64, GetPostTagsCallbackFunc) error
	GetPostTagsForPost(context.Context, int64, GetPostTagsCallbackFunc) error
	GetPostTagWithId(context.Context, int64) (*activitypub.PostTag, error)
	AddPostTag(context.Context, *activitypub.PostTag) error
	RemovePostTag(context.Context, *activitypub.PostTag) error
	GetPostTagsAll(context.Context, GetPostTagsCallbackFunc) error
	Close(context.Context) error
}

var post_tags_database_roster roster.Roster

// PostTagsDatabaseInitializationFunc is a function defined by individual post_tags_database package and used to create
// an instance of that post_tags_database
type PostTagsDatabaseInitializationFunc func(ctx context.Context, uri string) (PostTagsDatabase, error)

// RegisterPostTagsDatabase registers 'scheme' as a key pointing to 'init_func' in an internal lookup table
// used to create new `PostTagsDatabase` instances by the `NewPostTagsDatabase` method.
func RegisterPostTagsDatabase(ctx context.Context, scheme string, init_func PostTagsDatabaseInitializationFunc) error {

	err := ensurePostTagsDatabaseRoster()

	if err != nil {
		return err
	}

	return post_tags_database_roster.Register(ctx, scheme, init_func)
}

func ensurePostTagsDatabaseRoster() error {

	if post_tags_database_roster == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		post_tags_database_roster = r
	}

	return nil
}

// NewPostTagsDatabase returns a new `PostTagsDatabase` instance configured by 'uri'. The value of 'uri' is parsed
// as a `url.URL` and its scheme is used as the key for a corresponding `PostTagsDatabaseInitializationFunc`
// function used to instantiate the new `PostTagsDatabase`. It is assumed that the scheme (and initialization
// function) have been registered by the `RegisterPostTagsDatabase` method.
func NewPostTagsDatabase(ctx context.Context, uri string) (PostTagsDatabase, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := post_tags_database_roster.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	init_func := i.(PostTagsDatabaseInitializationFunc)
	return init_func(ctx, uri)
}

// Schemes returns the list of schemes that have been registered.
func PostTagsDatabaseSchemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensurePostTagsDatabaseRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range post_tags_database_roster.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}

func MigratePostTagsDatabaseFromURIs(ctx context.Context, from_uri string, to_uri string, count *int64, success *int64, errors *int64) error {

	from_ctx, from_cancel := context.WithTimeout(ctx, 5*time.Second)
	defer from_cancel()

	from_db, err := NewPostTagsDatabase(from_ctx, from_uri)

	if err != nil {
		return fmt.Errorf("Failed to create from database, %w", err)
	}

	defer from_db.Close(ctx)

	slog.Debug("Set up to database")

	to_ctx, to_cancel := context.WithTimeout(ctx, 5*time.Second)
	defer to_cancel()

	to_db, err := NewPostTagsDatabase(to_ctx, to_uri)

	if err != nil {
		return fmt.Errorf("Failed to create to database, %w", err)
	}

	defer to_db.Close(ctx)

	return MigratePostTagsDatabase(ctx, from_db, to_db, count, success, errors)
}

func MigratePostTagsDatabase(ctx context.Context, from_db PostTagsDatabase, to_db PostTagsDatabase, count *int64, success *int64, errors *int64) error {

	cb := func(ctx context.Context, t *activitypub.PostTag) error {

		defer atomic.AddInt64(count, 1)

		slog.Debug("Add", "tag", t.Id)
		err := to_db.AddPostTag(ctx, t)

		if err != nil {
			slog.Error("Failed to add tag", "tag", t.Id, "error", err)
			atomic.AddInt64(errors, 1)
		} else {
			atomic.AddInt64(success, 1)
		}

		return nil
	}

	slog.Debug("Retrieve tags")
	return from_db.GetPostTagsAll(ctx, cb)
}
