package activitypub

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/aaronland/go-roster"
)

type GetPostTagsCallbackFunc func(context.Context, *PostTag) error

type PostTagsDatabase interface {
	GetPostTagsForName(context.Context, string, GetPostTagsCallbackFunc) error
	GetPostTagsForAccount(context.Context, int64, GetPostTagsCallbackFunc) error
	GetPostTagsForPost(context.Context, int64, GetPostTagsCallbackFunc) error
	GetPostTagWithId(context.Context, int64) (*PostTag, error)
	AddPostTag(context.Context, *PostTag) error
	RemovePostTag(context.Context, *PostTag) error
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
