package database

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"sort"
	"strings"
	"sync/atomic"
	"time"

	"github.com/aaronland/go-roster"
	"github.com/sfomuseum/go-activitypub"
)

type GetMessageIdsCallbackFunc func(context.Context, int64) error
type GetMessagesCallbackFunc func(context.Context, *activitypub.Message) error

type MessagesDatabase interface {
	GetMessagesAll(context.Context, GetMessagesCallbackFunc) error
	GetMessageIdsForDateRange(context.Context, int64, int64, GetMessageIdsCallbackFunc) error
	GetMessagesForAccount(context.Context, int64, GetMessagesCallbackFunc) error
	GetMessagesForAccountAndAuthor(context.Context, int64, string, GetMessagesCallbackFunc) error
	GetMessageWithId(context.Context, int64) (*activitypub.Message, error)
	GetMessageWithAccountAndNoteIds(context.Context, int64, int64) (*activitypub.Message, error)
	AddMessage(context.Context, *activitypub.Message) error
	UpdateMessage(context.Context, *activitypub.Message) error
	RemoveMessage(context.Context, *activitypub.Message) error
	Close(context.Context) error
}

var messages_database_roster roster.Roster

// MessagesDatabaseInitializationFunc is a function defined by individual messages_database package and used to create
// an instance of that messages_database
type MessagesDatabaseInitializationFunc func(ctx context.Context, uri string) (MessagesDatabase, error)

// RegisterMessagesDatabase registers 'scheme' as a key pointing to 'init_func' in an internal lookup table
// used to create new `MessagesDatabase` instances by the `NewMessagesDatabase` method.
func RegisterMessagesDatabase(ctx context.Context, scheme string, init_func MessagesDatabaseInitializationFunc) error {

	err := ensureMessagesDatabaseRoster()

	if err != nil {
		return err
	}

	return messages_database_roster.Register(ctx, scheme, init_func)
}

func ensureMessagesDatabaseRoster() error {

	if messages_database_roster == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		messages_database_roster = r
	}

	return nil
}

// NewMessagesDatabase returns a new `MessagesDatabase` instance configured by 'uri'. The value of 'uri' is parsed
// as a `url.URL` and its scheme is used as the key for a corresponding `MessagesDatabaseInitializationFunc`
// function used to instantiate the new `MessagesDatabase`. It is assumed that the scheme (and initialization
// function) have been registered by the `RegisterMessagesDatabase` method.
func NewMessagesDatabase(ctx context.Context, uri string) (MessagesDatabase, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := messages_database_roster.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	init_func := i.(MessagesDatabaseInitializationFunc)
	return init_func(ctx, uri)
}

// Schemes returns the list of schemes that have been registered.
func MessagesDatabaseSchemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureMessagesDatabaseRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range messages_database_roster.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}

func MigrateMessagesDatabaseFromURIs(ctx context.Context, from_uri string, to_uri string, count *int64, success *int64, errors *int64) error {

	from_ctx, from_cancel := context.WithTimeout(ctx, 5*time.Second)
	defer from_cancel()

	from_db, err := NewMessagesDatabase(from_ctx, from_uri)

	if err != nil {
		return fmt.Errorf("Failed to create from database, %w", err)
	}

	defer from_db.Close(ctx)

	slog.Debug("Set up to database")

	to_ctx, to_cancel := context.WithTimeout(ctx, 5*time.Second)
	defer to_cancel()

	to_db, err := NewMessagesDatabase(to_ctx, to_uri)

	if err != nil {
		return fmt.Errorf("Failed to create to database, %w", err)
	}

	defer to_db.Close(ctx)

	return MigrateMessagesDatabase(ctx, from_db, to_db, count, success, errors)
}

func MigrateMessagesDatabase(ctx context.Context, from_db MessagesDatabase, to_db MessagesDatabase, count *int64, success *int64, errors *int64) error {

	cb := func(ctx context.Context, l *activitypub.Message) error {

		defer atomic.AddInt64(count, 1)

		slog.Debug("Add", "message", l.Id)
		err := to_db.AddMessage(ctx, l)

		if err != nil {
			slog.Error("Failed to add messages", "message", l.Id, "error", err)
			atomic.AddInt64(errors, 1)
		} else {
			atomic.AddInt64(success, 1)
		}

		return nil
	}

	slog.Debug("Retrieve messages")
	return from_db.GetMessagesAll(ctx, cb)
}
