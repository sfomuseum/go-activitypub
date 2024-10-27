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

// GetAccountIdsCallbackFunc is a custom function used to process an `activitypub.Account` identifier.
type GetAccountIdsCallbackFunc func(context.Context, int64) error

// GetAccountsCallbackFunc is a custom function used to process an `activitypub.Account` record.
type GetAccountsCallbackFunc func(context.Context, *activitypub.Account) error

// AccountsDatabase defines an interface for working with individual accounts associated with an atomic instance of the `go-activity` tools and services.
type AccountsDatabase interface {
	// GetAccounts iterates through all the account records and dispatches each to an instance of `GetAccountsCallbackFunc`.
	GetAccounts(context.Context, GetAccountsCallbackFunc) error
	// GetAccountsForDateRange iterates through all the account records created between two dates and dispatches each to an instance of `GetAccountIdsCallbackFunc`.
	GetAccountIdsForDateRange(context.Context, int64, int64, GetAccountIdsCallbackFunc) error
	// GetAccountWithId returns the account matching a specific 64-bit ID.
	GetAccountWithId(context.Context, int64) (*activitypub.Account, error)
	// GetAccountWithId returns the account matching a specific name.
	GetAccountWithName(context.Context, string) (*activitypub.Account, error)
	// AddAccount adds a new `activitypub.Account` instance.
	AddAccount(context.Context, *activitypub.Account) error
	// RemoveAccount removes a specific `activitypub.Account` instance.
	RemoveAccount(context.Context, *activitypub.Account) error
	// UpdateAccount updates a specific `activitypub.Account` instance.
	UpdateAccount(context.Context, *activitypub.Account) error
	// Close performs any final operations to terminate the underlying database connection.
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
