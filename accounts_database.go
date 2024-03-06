package activitypub

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/aaronland/go-roster"
)

type GetAccountIdsCallbackFunc func(context.Context, int64) error

type AccountsDatabase interface {
	GetAccountIdsForDateRange(context.Context, int64, int64, GetAccountIdsCallbackFunc) error
	GetAccountWithId(context.Context, int64) (*Account, error)
	GetAccountWithName(context.Context, string) (*Account, error)
	AddAccount(context.Context, *Account) error
	RemoveAccount(context.Context, *Account) error
	UpdateAccount(context.Context, *Account) error
	Close(context.Context) error
}

var account_database_roster roster.Roster

// AccountsDatabaseInitializationFunc is a function defined by individual account_database package and used to create
// an instance of that account_database
type AccountsDatabaseInitializationFunc func(ctx context.Context, uri string) (AccountsDatabase, error)

// RegisterAccountsDatabase registers 'scheme' as a key pointing to 'init_func' in an internal lookup table
// used to create new `AccountsDatabase` instances by the `NewAccountsDatabase` method.
func RegisterAccountsDatabase(ctx context.Context, scheme string, init_func AccountsDatabaseInitializationFunc) error {

	err := ensureAccountsDatabaseRoster()

	if err != nil {
		return err
	}

	return account_database_roster.Register(ctx, scheme, init_func)
}

func ensureAccountsDatabaseRoster() error {

	if account_database_roster == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		account_database_roster = r
	}

	return nil
}

// NewAccountsDatabase returns a new `AccountsDatabase` instance configured by 'uri'. The value of 'uri' is parsed
// as a `url.URL` and its scheme is used as the key for a corresponding `AccountsDatabaseInitializationFunc`
// function used to instantiate the new `AccountsDatabase`. It is assumed that the scheme (and initialization
// function) have been registered by the `RegisterAccountsDatabase` method.
func NewAccountsDatabase(ctx context.Context, uri string) (AccountsDatabase, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := account_database_roster.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	init_func := i.(AccountsDatabaseInitializationFunc)
	return init_func(ctx, uri)
}

// Schemes returns the list of schemes that have been registered.
func AccountsDatabaseSchemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureAccountsDatabaseRoster()

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
