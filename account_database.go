package activitypub

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/aaronland/go-roster"
)

type AccountDatabase interface {
	GetAccount(context.Context, string) (*Account, error)
	AddAccount(context.Context, *Account) error
	RemoveAccount(context.Context, *Account) error
	UpdateAccount(context.Context, *Account) error
}

var account_database_roster roster.Roster

// AccountDatabaseInitializationFunc is a function defined by individual account_database package and used to create
// an instance of that account_database
type AccountDatabaseInitializationFunc func(ctx context.Context, uri string) (AccountDatabase, error)

// RegisterAccountDatabase registers 'scheme' as a key pointing to 'init_func' in an internal lookup table
// used to create new `AccountDatabase` instances by the `NewAccountDatabase` method.
func RegisterAccountDatabase(ctx context.Context, scheme string, init_func AccountDatabaseInitializationFunc) error {

	err := ensureAccountDatabaseRoster()

	if err != nil {
		return err
	}

	return account_database_roster.Register(ctx, scheme, init_func)
}

func ensureAccountDatabaseRoster() error {

	if account_database_roster == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		account_database_roster = r
	}

	return nil
}

// NewAccountDatabase returns a new `AccountDatabase` instance configured by 'uri'. The value of 'uri' is parsed
// as a `url.URL` and its scheme is used as the key for a corresponding `AccountDatabaseInitializationFunc`
// function used to instantiate the new `AccountDatabase`. It is assumed that the scheme (and initialization
// function) have been registered by the `RegisterAccountDatabase` method.
func NewAccountDatabase(ctx context.Context, uri string) (AccountDatabase, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := account_database_roster.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	init_func := i.(AccountDatabaseInitializationFunc)
	return init_func(ctx, uri)
}

// Schemes returns the list of schemes that have been registered.
func AccountDatabaseSchemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureAccountDatabaseRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range account_database_roster.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}
