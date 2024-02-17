package activitypub

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"

	"github.com/sfomuseum/go-activitypub/sqlite"
)

const SQL_NOTES_TABLE_NAME string = "notes"

type SQLNotesDatabase struct {
	NotesDatabase
	database *sql.DB
}

func init() {
	ctx := context.Background()
	RegisterNotesDatabase(ctx, "sql", NewSQLNotesDatabase)
}

func NewSQLNotesDatabase(ctx context.Context, uri string) (NotesDatabase, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	engine := u.Host

	q := u.Query()
	dsn := q.Get("dsn")

	conn, err := sql.Open(engine, dsn)

	if err != nil {
		return nil, fmt.Errorf("Failed to open database connection, %w", err)
	}

	if engine == "sqlite3" {

		err := sqlite.SetupConnection(ctx, conn)

		if err != nil {
			return nil, fmt.Errorf("Failed to live hard and die fast, %w", err)
		}
	}

	db := &SQLNotesDatabase{
		database: conn,
	}

	return db, nil
}

func (db *SQLNotesDatabase) GetNoteWithId(ctx context.Context, note_id int64) (*Note, error) {
	where := "id = ?"
	return db.getNote(ctx, where, note_id)
}

func (db *SQLNotesDatabase) GetNoteWithNoteIdAndAuthor(ctx context.Context, note_id string, author_id string) (*Note, error) {

	where := "id = ?"
	return db.getNote(ctx, where, note_id, author_id)

}

func (db *SQLNotesDatabase) getNote(ctx context.Context, where string, args ...interface{}) (*Note, error) {

	q := fmt.Sprintf("SELECT id, note_id, author_id, body, created, lastmodified FROM %s WHERE %s", SQL_NOTES_TABLE_NAME, where)
	row := db.database.QueryRowContext(ctx, q, args...)

	var id int64
	var note_id string
	var author_id string
	var body []byte
	var created int64
	var lastmod int64

	err := row.Scan(&id, &note_id, &author_id, &body, &created, &lastmod)

	switch {
	case err == sql.ErrNoRows:
		return nil, ErrNotFound
	case err != nil:
		return nil, fmt.Errorf("Failed to query database, %w", err)
	default:

		n := &Note{
			Id:           id,
			NoteId:       note_id,
			AuthorId:     author_id,
			Body:         body,
			Created:      created,
			LastModified: lastmod,
		}

		return n, nil
	}

}

func (db *SQLNotesDatabase) AddNote(ctx context.Context, note *Note) error {

	q := fmt.Sprintf("INSERT INTO %s (id, note_id, author_id, body, created, lastmodified) VALUES (?, ?, ?, ?, ?, ?)", SQL_NOTES_TABLE_NAME)

	_, err := db.database.ExecContext(ctx, q, note.Id, note.NoteId, note.AuthorId, note.Body, note.Created, note.LastModified)

	if err != nil {
		return fmt.Errorf("Failed to add note, %w", err)
	}

	return nil
}

func (db *SQLNotesDatabase) UpdateNote(ctx context.Context, note *Note) error {

	q := fmt.Sprintf("UPDATE %s SET note_id=?, author_id=?, body=?, created=?, lastmodified=? WHERE id = ?", SQL_NOTES_TABLE_NAME)

	_, err := db.database.ExecContext(ctx, q, note.NoteId, note.AuthorId, note.Body, note.Created, note.LastModified, note.Id)

	if err != nil {
		return fmt.Errorf("Failed to add note, %w", err)
	}

	return nil
}

func (db *SQLNotesDatabase) RemoveNote(ctx context.Context, note *Note) error {

	q := fmt.Sprintf("DELETE FROM %s WHERE id = ?", SQL_NOTES_TABLE_NAME)

	_, err := db.database.ExecContext(ctx, q, note.Id)

	if err != nil {
		return fmt.Errorf("Failed to remove note, %w", err)
	}

	return nil
}
