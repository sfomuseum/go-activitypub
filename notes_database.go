package activitypub

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/aaronland/go-roster"
)

type GetNoteIdsCallbackFunc func(context.Context, int64) error
type GetNotesCallbackFunc func(context.Context, *Note) error

type NotesDatabase interface {
	GetNoteIdsForDateRange(context.Context, int64, int64, GetNoteIdsCallbackFunc) error
	GetNoteWithId(context.Context, int64) (*Note, error)
	GetNoteWithUUIDAndAuthorAddress(context.Context, string, string) (*Note, error)
	AddNote(context.Context, *Note) error
	UpdateNote(context.Context, *Note) error
	RemoveNote(context.Context, *Note) error
	Close(context.Context) error
}

var notes_database_roster roster.Roster

// NotesDatabaseInitializationFunc is a function defined by individual notes_database package and used to create
// an instance of that notes_database
type NotesDatabaseInitializationFunc func(ctx context.Context, uri string) (NotesDatabase, error)

// RegisterNotesDatabase registers 'scheme' as a key pointing to 'init_func' in an internal lookup table
// used to create new `NotesDatabase` instances by the `NewNotesDatabase` method.
func RegisterNotesDatabase(ctx context.Context, scheme string, init_func NotesDatabaseInitializationFunc) error {

	err := ensureNotesDatabaseRoster()

	if err != nil {
		return err
	}

	return notes_database_roster.Register(ctx, scheme, init_func)
}

func ensureNotesDatabaseRoster() error {

	if notes_database_roster == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		notes_database_roster = r
	}

	return nil
}

// NewNotesDatabase returns a new `NotesDatabase` instance configured by 'uri'. The value of 'uri' is parsed
// as a `url.URL` and its scheme is used as the key for a corresponding `NotesDatabaseInitializationFunc`
// function used to instantiate the new `NotesDatabase`. It is assumed that the scheme (and initialization
// function) have been registered by the `RegisterNotesDatabase` method.
func NewNotesDatabase(ctx context.Context, uri string) (NotesDatabase, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := notes_database_roster.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	init_func := i.(NotesDatabaseInitializationFunc)
	return init_func(ctx, uri)
}

// Schemes returns the list of schemes that have been registered.
func NotesDatabaseSchemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureNotesDatabaseRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range notes_database_roster.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}
