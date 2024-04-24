package activitypub

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/aaronland/go-roster"
)

type GetPropertiesCallbackFunc func(context.Context, *Property) error

type PropertiesDatabase interface {
	GetPropertiesForAccount(context.Context, int64, GetPropertiesCallbackFunc) error
	AddProperty(context.Context, *Property) error
	UpdateProperty(context.Context, *Property) error
	RemoveProperty(context.Context, *Property) error
	Close(context.Context) error
}

var properties_database_roster roster.Roster

// PropertiesDatabaseInitializationFunc is a function defined by individual properties_database package and used to create
// an instance of that properties_database
type PropertiesDatabaseInitializationFunc func(ctx context.Context, uri string) (PropertiesDatabase, error)

// RegisterPropertiesDatabase registers 'scheme' as a key pointing to 'init_func' in an internal lookup table
// used to create new `PropertiesDatabase` instances by the `NewPropertiesDatabase` method.
func RegisterPropertiesDatabase(ctx context.Context, scheme string, init_func PropertiesDatabaseInitializationFunc) error {

	err := ensurePropertiesDatabaseRoster()

	if err != nil {
		return err
	}

	return properties_database_roster.Register(ctx, scheme, init_func)
}

func ensurePropertiesDatabaseRoster() error {

	if properties_database_roster == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		properties_database_roster = r
	}

	return nil
}

// NewPropertiesDatabase returns a new `PropertiesDatabase` instance configured by 'uri'. The value of 'uri' is parsed
// as a `url.URL` and its scheme is used as the key for a corresponding `PropertiesDatabaseInitializationFunc`
// function used to instantiate the new `PropertiesDatabase`. It is assumed that the scheme (and initialization
// function) have been registered by the `RegisterPropertiesDatabase` method.
func NewPropertiesDatabase(ctx context.Context, uri string) (PropertiesDatabase, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := properties_database_roster.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	init_func := i.(PropertiesDatabaseInitializationFunc)
	return init_func(ctx, uri)
}

// Schemes returns the list of schemes that have been registered.
func PropertiesDatabaseSchemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensurePropertiesDatabaseRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range properties_database_roster.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}
