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

type ActivitiesDatabase interface {
	AddRecord(context.Context, *activitypub.Activity) error
	RemoveRecord(context.Context, *activitypub.Activity) error
	UpdateRecord(context.Context, *activitypub.Activity) error
	GetRecord(context.Context, int64) (*activitypub.Activity, error)
	QueryRecords(context.Context, *Query) iter.Seq2[*activitypub.Activity, error]
	Close() error

	GetActivityWithActivityPubId(context.Context, string) (*activitypub.Activity, error)
	GetActivityWithActivityTypeAnId(context.Context, activitypub.ActivityType, int64) (*activitypub.Activity, error)
	GetActivitiesForAccount(context.Context, int64) iter.Seq2[*activitypub.Activity, error]
}

var activities_database_roster roster.Roster

// ActivitiesDatabaseInitializationFunc is a function defined by individual activities_database package and used to create
// an instance of that activities_database
type ActivitiesDatabaseInitializationFunc func(ctx context.Context, uri string) (ActivitiesDatabase, error)

// RegisterActivitiesDatabase registers 'scheme' as a key pointing to 'init_func' in an internal lookup table
// used to create new `ActivitiesDatabase` instances by the `NewActivitiesDatabase` method.
func RegisterActivitiesDatabase(ctx context.Context, scheme string, init_func ActivitiesDatabaseInitializationFunc) error {

	err := ensureActivitiesDatabaseRoster()

	if err != nil {
		return err
	}

	return activities_database_roster.Register(ctx, scheme, init_func)
}

func ensureActivitiesDatabaseRoster() error {

	if activities_database_roster == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		activities_database_roster = r
	}

	return nil
}

// NewActivitiesDatabase returns a new `ActivitiesDatabase` instance configured by 'uri'. The value of 'uri' is parsed
// as a `url.URL` and its scheme is used as the key for a corresponding `ActivitiesDatabaseInitializationFunc`
// function used to instantiate the new `ActivitiesDatabase`. It is assumed that the scheme (and initialization
// function) have been registered by the `RegisterActivitiesDatabase` method.
func NewActivitiesDatabase(ctx context.Context, uri string) (ActivitiesDatabase, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := activities_database_roster.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	init_func := i.(ActivitiesDatabaseInitializationFunc)
	return init_func(ctx, uri)
}

// Schemes returns the list of schemes that have been registered.
func ActivitiesDatabaseSchemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureActivitiesDatabaseRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range activities_database_roster.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}

func MigrateActivitiesDatabaseFromURIs(ctx context.Context, from_uri string, to_uri string, count *int64, success *int64, errors *int64) error {

	from_ctx, from_cancel := context.WithTimeout(ctx, 5*time.Second)
	defer from_cancel()

	from_db, err := NewActivitiesDatabase(from_ctx, from_uri)

	if err != nil {
		return fmt.Errorf("Failed to create from database, %w", err)
	}

	defer from_db.Close()

	slog.Debug("Set up to database")

	to_ctx, to_cancel := context.WithTimeout(ctx, 5*time.Second)
	defer to_cancel()

	to_db, err := NewActivitiesDatabase(to_ctx, to_uri)

	if err != nil {
		return fmt.Errorf("Failed to create to database, %w", err)
	}

	defer to_db.Close()

	_, _, _, err = Migrate(ctx, from_db, to_db)
	return err
}
