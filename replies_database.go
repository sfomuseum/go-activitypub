package activitypub

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/aaronland/go-roster"
)

type GetRepliesCallbackFunc func(context.Context, *Reply) error

type RepliesDatabase interface {
	GetRepliesForActor(context.Context, string, GetRepliesCallbackFunc) error
	GetRepliesForAccount(context.Context, int64, GetRepliesCallbackFunc) error
	GetRepliesForPost(context.Context, int64, GetRepliesCallbackFunc) error
	GetReplyWithId(context.Context, int64) (*Reply, error)
	GetReplyWithReplyId(context.Context, string) (*Reply, error)
	AddReply(context.Context, *Reply) error
	RemoveReply(context.Context, *Reply) error
	Close(context.Context) error
}

var reply_database_roster roster.Roster

// RepliesDatabaseInitializationFunc is a function defined by individual reply_database package and used to create
// an instance of that reply_database
type RepliesDatabaseInitializationFunc func(ctx context.Context, uri string) (RepliesDatabase, error)

// RegisterRepliesDatabase registers 'scheme' as a key pointing to 'init_func' in an internal lookup table
// used to create new `RepliesDatabase` instances by the `NewRepliesDatabase` method.
func RegisterRepliesDatabase(ctx context.Context, scheme string, init_func RepliesDatabaseInitializationFunc) error {

	err := ensureRepliesDatabaseRoster()

	if err != nil {
		return err
	}

	return reply_database_roster.Register(ctx, scheme, init_func)
}

func ensureRepliesDatabaseRoster() error {

	if reply_database_roster == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		reply_database_roster = r
	}

	return nil
}

// NewRepliesDatabase returns a new `RepliesDatabase` instance configured by 'uri'. The value of 'uri' is parsed
// as a `url.URL` and its scheme is used as the key for a corresponding `RepliesDatabaseInitializationFunc`
// function used to instantiate the new `RepliesDatabase`. It is assumed that the scheme (and initialization
// function) have been registered by the `RegisterRepliesDatabase` method.
func NewRepliesDatabase(ctx context.Context, uri string) (RepliesDatabase, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := reply_database_roster.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	init_func := i.(RepliesDatabaseInitializationFunc)
	return init_func(ctx, uri)
}

// Schemes returns the list of schemes that have been registered.
func RepliesDatabaseSchemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureRepliesDatabaseRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range reply_database_roster.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}
