package database

import (
	"context"
	"fmt"
	"io"

	aa_docstore "github.com/aaronland/gocloud-docstore"
	"github.com/sfomuseum/go-activitypub"
	gc_docstore "gocloud.dev/docstore"
)

type DocstoreNotesDatabase struct {
	NotesDatabase
	collection *gc_docstore.Collection
}

func init() {

	ctx := context.Background()

	err := RegisterNotesDatabase(ctx, "awsdynamodb", NewDocstoreNotesDatabase)

	if err != nil {
		panic(err)
	}

	for _, scheme := range gc_docstore.DefaultURLMux().CollectionSchemes() {
		err := RegisterNotesDatabase(ctx, scheme, NewDocstoreNotesDatabase)

		if err != nil {
			panic(err)
		}
	}
}

func NewDocstoreNotesDatabase(ctx context.Context, uri string) (NotesDatabase, error) {

	col, err := aa_docstore.OpenCollection(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to open collection, %w", err)
	}

	db := &DocstoreNotesDatabase{
		collection: col,
	}

	return db, nil
}

func (db *DocstoreNotesDatabase) GetNoteIdsForDateRange(ctx context.Context, start int64, end int64, cb GetNoteIdsCallbackFunc) error {

	q := db.collection.Query()
	q = q.Where("Created", ">=", start)
	q = q.Where("Created", "<=", end)

	iter := q.Get(ctx, "Id")
	defer iter.Stop()

	for {

		var n activitypub.Note
		err := iter.Next(ctx, &n)

		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("Failed to interate, %w", err)
		} else {
			err := cb(ctx, n.Id)

			if err != nil {
				return fmt.Errorf("Failed to invoke callback for note %d, %w", n.Id, err)
			}
		}
	}

	return nil
}

func (db *DocstoreNotesDatabase) GetNoteWithId(ctx context.Context, note_id int64) (*activitypub.Note, error) {

	q := db.collection.Query()
	q = q.Where("Id", "=", note_id)

	return db.getNote(ctx, q)
}

func (db *DocstoreNotesDatabase) GetNoteWithUUIDAndAuthorAddress(ctx context.Context, note_uuid string, author_address string) (*activitypub.Note, error) {

	q := db.collection.Query()
	q = q.Where("UUID", "=", note_uuid)
	q = q.Where("AuthorAddress", "=", author_address)

	return db.getNote(ctx, q)
}

func (db *DocstoreNotesDatabase) getNote(ctx context.Context, q *gc_docstore.Query) (*activitypub.Note, error) {

	iter := q.Get(ctx)
	defer iter.Stop()

	for {

		var n activitypub.Note
		err := iter.Next(ctx, &n)

		if err == io.EOF {
			break
		} else if err != nil {
			return nil, fmt.Errorf("Failed to interate, %w", err)
		} else {
			return &n, nil
		}
	}

	return nil, activitypub.ErrNotFound
}

func (db *DocstoreNotesDatabase) AddNote(ctx context.Context, note *activitypub.Note) error {

	return db.collection.Put(ctx, note)
}

func (db *DocstoreNotesDatabase) UpdateNote(ctx context.Context, note *activitypub.Note) error {

	return db.collection.Replace(ctx, note)
}

func (db *DocstoreNotesDatabase) RemoveNote(ctx context.Context, note *activitypub.Note) error {

	return db.collection.Delete(ctx, note)
}

func (db *DocstoreNotesDatabase) Close(ctx context.Context) error {
	return db.collection.Close()
}
