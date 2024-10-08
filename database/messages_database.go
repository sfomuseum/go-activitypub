package activitypub

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/aaronland/go-roster"
)

type GetMessageIdsCallbackFunc func(context.Context, int64) error
type GetMessagesCallbackFunc func(context.Context, *Message) error

type MessagesDatabase interface {
	GetMessageIdsForDateRange(context.Context, int64, int64, GetMessageIdsCallbackFunc) error
	GetMessagesForAccount(context.Context, int64, GetMessagesCallbackFunc) error
	GetMessagesForAccountAndAuthor(context.Context, int64, string, GetMessagesCallbackFunc) error
	GetMessageWithId(context.Context, int64) (*Message, error)
	GetMessageWithAccountAndNoteIds(context.Context, int64, int64) (*Message, error)
	AddMessage(context.Context, *Message) error
	UpdateMessage(context.Context, *Message) error
	RemoveMessage(context.Context, *Message) error
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
