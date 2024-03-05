package activitypub

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/aaronland/go-roster"
)

type PostsDatabase interface {
	GetPostWithId(context.Context, int64) (*Post, error)
	AddPost(context.Context, *Post) error
	RemovePost(context.Context, *Post) error
	UpdatePost(context.Context, *Post) error
	Close(context.Context) error
}

var post_database_roster roster.Roster

// PostsDatabaseInitializationFunc is a function defined by individual post_database package and used to create
// an instance of that post_database
type PostsDatabaseInitializationFunc func(ctx context.Context, uri string) (PostsDatabase, error)

// RegisterPostsDatabase registers 'scheme' as a key pointing to 'init_func' in an internal lookup table
// used to create new `PostsDatabase` instances by the `NewPostsDatabase` method.
func RegisterPostsDatabase(ctx context.Context, scheme string, init_func PostsDatabaseInitializationFunc) error {

	err := ensurePostsDatabaseRoster()

	if err != nil {
		return err
	}

	return post_database_roster.Register(ctx, scheme, init_func)
}

func ensurePostsDatabaseRoster() error {

	if post_database_roster == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		post_database_roster = r
	}

	return nil
}

// NewPostsDatabase returns a new `PostsDatabase` instance configured by 'uri'. The value of 'uri' is parsed
// as a `url.URL` and its scheme is used as the key for a corresponding `PostsDatabaseInitializationFunc`
// function used to instantiate the new `PostsDatabase`. It is assumed that the scheme (and initialization
// function) have been registered by the `RegisterPostsDatabase` method.
func NewPostsDatabase(ctx context.Context, uri string) (PostsDatabase, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := post_database_roster.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	init_func := i.(PostsDatabaseInitializationFunc)
	return init_func(ctx, uri)
}

// Schemes returns the list of schemes that have been registered.
func PostsDatabaseSchemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensurePostsDatabaseRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range post_database_roster.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}
