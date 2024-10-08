package activitypub

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/aaronland/go-roster"
)

type GetAliasesCallbackFunc func(context.Context, *Alias) error

type AliasesDatabase interface {
	GetAliasesForAccount(context.Context, int64, GetAliasesCallbackFunc) error
	GetAliasWithName(context.Context, string) (*Alias, error)
	AddAlias(context.Context, *Alias) error
	RemoveAlias(context.Context, *Alias) error
	Close(context.Context) error
}

var alias_database_roster roster.Roster

// AliasesDatabaseInitializationFunc is a function defined by individual alias_database package and used to create
// an instance of that alias_database
type AliasesDatabaseInitializationFunc func(ctx context.Context, uri string) (AliasesDatabase, error)

// RegisterAliasesDatabase registers 'scheme' as a key pointing to 'init_func' in an internal lookup table
// used to create new `AliasesDatabase` instances by the `NewAliasesDatabase` method.
func RegisterAliasesDatabase(ctx context.Context, scheme string, init_func AliasesDatabaseInitializationFunc) error {

	err := ensureAliasesDatabaseRoster()

	if err != nil {
		return err
	}

	return alias_database_roster.Register(ctx, scheme, init_func)
}

func ensureAliasesDatabaseRoster() error {

	if alias_database_roster == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		alias_database_roster = r
	}

	return nil
}

// NewAliasesDatabase returns a new `AliasesDatabase` instance configured by 'uri'. The value of 'uri' is parsed
// as a `url.URL` and its scheme is used as the key for a corresponding `AliasesDatabaseInitializationFunc`
// function used to instantiate the new `AliasesDatabase`. It is assumed that the scheme (and initialization
// function) have been registered by the `RegisterAliasesDatabase` method.
func NewAliasesDatabase(ctx context.Context, uri string) (AliasesDatabase, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := alias_database_roster.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	init_func := i.(AliasesDatabaseInitializationFunc)
	return init_func(ctx, uri)
}

// Schemes returns the list of schemes that have been registered.
func AliasesDatabaseSchemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureAliasesDatabaseRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range alias_database_roster.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}
