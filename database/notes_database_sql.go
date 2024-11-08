package database

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"

	pg_sql "github.com/aaronland/go-pagination-sql"
	"github.com/aaronland/go-pagination/countable"
	"github.com/sfomuseum/go-activitypub"
	"github.com/sfomuseum/go-activitypub/sqlite"
)

const SQL_NOTES_TABLE_NAME string = "notes"

type SQLNotesDatabase struct {
	NotesDatabase
	database *sql.DB
}

func init() {
	ctx := context.Background()
	err := RegisterNotesDatabase(ctx, "sql", NewSQLNotesDatabase)

	if err != nil {
		panic(err)
	}
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

func (db *SQLNotesDatabase) GetNoteIdsForDateRange(ctx context.Context, start int64, end int64, cb GetNoteIdsCallbackFunc) error {

	pg_callback := func(pg_rsp pg_sql.PaginatedResponse) error {

		rows := pg_rsp.Rows()

		for rows.Next() {

			var id int64

			err := rows.Scan(&id)

			if err != nil {
				return fmt.Errorf("Failed to query database, %w", err)
			}

			err = cb(ctx, id)

			if err != nil {
				return fmt.Errorf("Failed to execute following callback for note %d, %w", id, err)
			}

			return nil
		}

		err := rows.Close()

		if err != nil {
			return fmt.Errorf("Failed to iterate through database rows, %w", err)
		}

		return nil
	}

	pg_opts, err := countable.NewCountableOptions()

	if err != nil {
		return fmt.Errorf("Failed to create pagination options, %w", err)
	}

	q := fmt.Sprintf("SELECT id FROM %s WHERE created >= ? AND created <= ?", SQL_NOTES_TABLE_NAME)

	err = pg_sql.QueryPaginatedAll(db.database, pg_opts, pg_callback, q, start, end)

	if err != nil {
		return fmt.Errorf("Failed to execute paginated query, %w", err)
	}

	return nil
}

func (db *SQLNotesDatabase) GetNoteWithId(ctx context.Context, note_id int64) (*activitypub.Note, error) {
	where := "id = ?"
	return db.getNote(ctx, where, note_id)
}

func (db *SQLNotesDatabase) GetNoteWithUUIDAndAuthorAddress(ctx context.Context, uuid string, author_address string) (*activitypub.Note, error) {

	// Note the order of arguments this is to account for the
	// notes_by_author_address index.

	where := "author_address=? AND uuid = ?"
	return db.getNote(ctx, where, author_address, uuid)
}

func (db *SQLNotesDatabase) getNote(ctx context.Context, where string, args ...interface{}) (*activitypub.Note, error) {

	q := fmt.Sprintf("SELECT id, uuid, author_address, body, created, lastmodified FROM %s WHERE %s", SQL_NOTES_TABLE_NAME, where)
	row := db.database.QueryRowContext(ctx, q, args...)

	var id int64
	var uuid string
	var author_address string
	var body string
	var created int64
	var lastmod int64

	err := row.Scan(&id, &uuid, &author_address, &body, &created, &lastmod)

	switch {
	case err == sql.ErrNoRows:
		return nil, activitypub.ErrNotFound
	case err != nil:
		return nil, fmt.Errorf("Failed to query database, %w", err)
	default:

		n := &activitypub.Note{
			Id:            id,
			UUID:          uuid,
			AuthorAddress: author_address,
			Body:          body,
			Created:       created,
			LastModified:  lastmod,
		}

		return n, nil
	}

}

func (db *SQLNotesDatabase) AddNote(ctx context.Context, note *activitypub.Note) error {

	q := fmt.Sprintf("INSERT INTO %s (id, uuid, author_address, body, created, lastmodified) VALUES (?, ?, ?, ?, ?, ?)", SQL_NOTES_TABLE_NAME)

	_, err := db.database.ExecContext(ctx, q, note.Id, note.UUID, note.AuthorAddress, note.Body, note.Created, note.LastModified)

	if err != nil {
		return fmt.Errorf("Failed to add note, %w", err)
	}

	return nil
}

func (db *SQLNotesDatabase) UpdateNote(ctx context.Context, note *activitypub.Note) error {

	q := fmt.Sprintf("UPDATE %s SET uuid=?, author_address=?, body=?, created=?, lastmodified=? WHERE id = ?", SQL_NOTES_TABLE_NAME)

	_, err := db.database.ExecContext(ctx, q, note.UUID, note.AuthorAddress, note.Body, note.Created, note.LastModified, note.Id)

	if err != nil {
		return fmt.Errorf("Failed to add note, %w", err)
	}

	return nil
}

func (db *SQLNotesDatabase) RemoveNote(ctx context.Context, note *activitypub.Note) error {

	q := fmt.Sprintf("DELETE FROM %s WHERE id = ?", SQL_NOTES_TABLE_NAME)

	_, err := db.database.ExecContext(ctx, q, note.Id)

	if err != nil {
		return fmt.Errorf("Failed to remove note, %w", err)
	}

	return nil
}

func (db *SQLNotesDatabase) Close(ctx context.Context) error {
	return db.database.Close()
}
