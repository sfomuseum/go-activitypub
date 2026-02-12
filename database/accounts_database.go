package database

import (
	"context"
	"fmt"
	"iter"
	"log/slog"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/aaronland/go-roster"
	"github.com/sfomuseum/go-activitypub"
)

// AccountsDatabase defines an interface for working with individual accounts associated with an atomic instance of the `go-activity` tools and services.
type AccountsDatabase interface {
	AddRecord(context.Context, *activitypub.Account) error
	RemoveRecord(context.Context, *activitypub.Account) error
	UpdateRecord(context.Context, *activitypub.Account) error
	GetRecord(context.Context, int64) (*activitypub.Account, error)
	QueryRecords(context.Context, *Query) iter.Seq2[*activitypub.Account, error]
	Close() error

	GetAccountIdsForDateRange(context.Context, int64, int64) iter.Seq2[int64, error]
	// GetAccountWithId returns the account matching a specific 64-bit ID.
	GetAccountWithId(context.Context, int64) (*activitypub.Account, error)
	// GetAccountWithId returns the account matching a specific name.
	GetAccountWithName(context.Context, string) (*activitypub.Account, error)
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

func MigrateAccountsDatabaseFromURIs(ctx context.Context, from_uri string, to_uri string, count *int64, success *int64, errors *int64) error {

	from_ctx, from_cancel := context.WithTimeout(ctx, 5*time.Second)
	defer from_cancel()

	from_db, err := NewAccountsDatabase(from_ctx, from_uri)

	if err != nil {
		return fmt.Errorf("Failed to create from database, %w", err)
	}

	defer from_db.Close()

	slog.Debug("Set up to database")

	to_ctx, to_cancel := context.WithTimeout(ctx, 5*time.Second)
	defer to_cancel()

	to_db, err := NewAccountsDatabase(to_ctx, to_uri)

	if err != nil {
		return fmt.Errorf("Failed to create to database, %w", err)
	}

	defer to_db.Close()

	_, _, _, err = Migrate(ctx, from_db, to_db)

	if err != nil {
		return err
	}

	return nil
}
