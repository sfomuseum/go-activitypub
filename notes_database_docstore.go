package activitypub

import (
	"context"
	"fmt"
	"io"

	aa_docstore "github.com/aaronland/gocloud-docstore"
	gc_docstore "gocloud.dev/docstore"
)

type DocstoreNotesDatabase struct {
	NotesDatabase
	collection *gc_docstore.Collection
}

func init() {

	ctx := context.Background()

	RegisterNotesDatabase(ctx, "awsdynamodb", NewDocstoreNotesDatabase)

	for _, scheme := range gc_docstore.DefaultURLMux().CollectionSchemes() {
		RegisterNotesDatabase(ctx, scheme, NewDocstoreNotesDatabase)
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

func (db *DocstoreNotesDatabase) GetNoteWithId(ctx context.Context, note_id int64) (*Note, error) {

	q := db.collection.Query()
	q = q.Where("Id", "=", note_id)

	return db.getNote(ctx, q)
}

func (db *DocstoreNotesDatabase) GetNoteWithUUIDAndAuthorAddress(ctx context.Context, note_uuid string, author_address string) (*Note, error) {

	q := db.collection.Query()
	q = q.Where("UUID", "=", note_uuid)
	q = q.Where("AuthorAddress", "=", author_address)

	return db.getNote(ctx, q)
}

func (db *DocstoreNotesDatabase) getNote(ctx context.Context, q *gc_docstore.Query) (*Note, error) {

	iter := q.Get(ctx)
	defer iter.Stop()

	for {

		var n Note
		err := iter.Next(ctx, &n)

		if err == io.EOF {
			break
		} else if err != nil {
			return nil, fmt.Errorf("Failed to interate, %w", err)
		} else {
			return &n, nil
		}
	}

	return nil, ErrNotFound
}

func (db *DocstoreNotesDatabase) AddNote(ctx context.Context, note *Note) error {

	return db.collection.Put(ctx, note)
}

func (db *DocstoreNotesDatabase) UpdateNote(ctx context.Context, note *Note) error {

	return db.collection.Replace(ctx, note)
}

func (db *DocstoreNotesDatabase) RemoveNote(ctx context.Context, note *Note) error {

	return db.collection.Delete(ctx, note)
}

func (db *DocstoreNotesDatabase) Close(ctx context.Context) error {
	return db.collection.Close()
}
