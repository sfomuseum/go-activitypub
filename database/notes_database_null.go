package database

import (
	"context"

	"github.com/sfomuseum/go-activitypub"
)

type NullNotesDatabase struct {
	NotesDatabase
}

func init() {

	ctx := context.Background()
	err := RegisterNotesDatabase(ctx, "null", NewNullNotesDatabase)

	if err != nil {
		panic(err)
	}
}

func NewNullNotesDatabase(ctx context.Context, uri string) (NotesDatabase, error) {

	db := &NullNotesDatabase{}
	return db, nil
}

func (db *NullNotesDatabase) GetNoteIdsForDateRange(ctx context.Context, start int64, end int64, cb GetNoteIdsCallbackFunc) error {
	return nil
}

func (db *NullNotesDatabase) GetNoteWithId(ctx context.Context, note_id int64) (*activitypub.Note, error) {
	return nil, activitypub.ErrNotFound
}

func (db *NullNotesDatabase) GetNoteWithUUIDAndAuthorAddress(ctx context.Context, note_uuid string, author_address string) (*activitypub.Note, error) {
	return nil, activitypub.ErrNotFound
}

func (db *NullNotesDatabase) AddNote(ctx context.Context, note *activitypub.Note) error {
	return nil
}

func (db *NullNotesDatabase) UpdateNote(ctx context.Context, note *activitypub.Note) error {
	return nil
}

func (db *NullNotesDatabase) RemoveNote(ctx context.Context, note *activitypub.Note) error {
	return nil
}

func (db *NullNotesDatabase) Close(ctx context.Context) error {
	return nil
}
