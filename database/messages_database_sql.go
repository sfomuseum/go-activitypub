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

const SQL_MESSAGES_TABLE_NAME string = "messages"

type SQLMessagesDatabase struct {
	MessagesDatabase
	database *sql.DB
}

func init() {
	ctx := context.Background()
	err := RegisterMessagesDatabase(ctx, "sql", NewSQLMessagesDatabase)

	if err != nil {
		panic(err)
	}
}

func NewSQLMessagesDatabase(ctx context.Context, uri string) (MessagesDatabase, error) {

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

	db := &SQLMessagesDatabase{
		database: conn,
	}

	return db, nil
}

func (db *SQLMessagesDatabase) GetMessageIdsForDateRange(ctx context.Context, start int64, end int64, cb GetMessageIdsCallbackFunc) error {

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
				return fmt.Errorf("Failed to execute following callback for message %d, %w", id, err)
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

	q := fmt.Sprintf("SELECT id FROM %s WHERE created >= ? AND created <= ?", SQL_MESSAGES_TABLE_NAME)

	err = pg_sql.QueryPaginatedAll(db.database, pg_opts, pg_callback, q, start, end)

	if err != nil {
		return fmt.Errorf("Failed to execute paginated query, %w", err)
	}

	return nil
}

func (db *SQLMessagesDatabase) GetMessageWithId(ctx context.Context, message_id int64) (*activitypub.Message, error) {
	where := "id = ?"
	return db.getMessage(ctx, where, message_id)
}

func (db *SQLMessagesDatabase) GetMessageWithAccountAndNoteIds(ctx context.Context, account_id int64, note_id int64) (*activitypub.Message, error) {

	where := "account_id = ? AND note_id = ?"
	return db.getMessage(ctx, where, account_id, note_id)

}

func (db *SQLMessagesDatabase) getMessage(ctx context.Context, where string, args ...interface{}) (*activitypub.Message, error) {

	q := fmt.Sprintf("SELECT id, note_id, author_address, account_id, created, lastmodified FROM %s WHERE %s", SQL_MESSAGES_TABLE_NAME, where)
	row := db.database.QueryRowContext(ctx, q, args...)

	var id int64
	var note_id int64
	var author_address string
	var account_id int64
	var created int64
	var lastmod int64

	err := row.Scan(&id, &note_id, &author_address, &account_id, &created, &lastmod)

	switch {
	case err == sql.ErrNoRows:
		return nil, activitypub.ErrNotFound
	case err != nil:
		return nil, fmt.Errorf("Failed to query database, %w", err)
	default:

		n := &activitypub.Message{
			Id:            id,
			NoteId:        note_id,
			AuthorAddress: author_address,
			AccountId:     account_id,
			Created:       created,
			LastModified:  lastmod,
		}

		return n, nil
	}

}

func (db *SQLMessagesDatabase) AddMessage(ctx context.Context, message *activitypub.Message) error {

	q := fmt.Sprintf("INSERT INTO %s (id, note_id, author_address, account_id, created, lastmodified) VALUES (?, ?, ?, ?, ?, ?)", SQL_MESSAGES_TABLE_NAME)

	_, err := db.database.ExecContext(ctx, q, message.Id, message.NoteId, message.AuthorAddress, message.AccountId, message.Created, message.LastModified)

	if err != nil {
		return fmt.Errorf("Failed to add message, %w", err)
	}

	return nil
}

func (db *SQLMessagesDatabase) UpdateMessage(ctx context.Context, message *activitypub.Message) error {

	q := fmt.Sprintf("UPDATE %s SET note_id=?, author_address=?, account_id=?, created=?, lastmodified=? WHERE id = ?", SQL_MESSAGES_TABLE_NAME)

	_, err := db.database.ExecContext(ctx, q, message.NoteId, message.AuthorAddress, message.AccountId, message.Created, message.LastModified, message.Id)

	if err != nil {
		return fmt.Errorf("Failed to add message, %w", err)
	}

	return nil
}

func (db *SQLMessagesDatabase) RemoveMessage(ctx context.Context, message *activitypub.Message) error {

	q := fmt.Sprintf("DELETE FROM %s WHERE id = ?", SQL_MESSAGES_TABLE_NAME)

	_, err := db.database.ExecContext(ctx, q, message.Id)

	if err != nil {
		return fmt.Errorf("Failed to remove message, %w", err)
	}

	return nil
}

func (db *SQLMessagesDatabase) GetMessagesForAccount(ctx context.Context, account_id int64, callback_func GetMessagesCallbackFunc) error {

	where := "account_id = ?"
	args := []interface{}{
		account_id,
	}

	return db.getMessagesWithCallback(ctx, where, args, callback_func)
}

func (db *SQLMessagesDatabase) GetMessagesForAccountAndAuthor(ctx context.Context, account_id int64, author_address string, callback_func GetMessagesCallbackFunc) error {

	where := "account_id = ? AND author_address = ?"
	args := []interface{}{
		account_id,
		author_address,
	}

	return db.getMessagesWithCallback(ctx, where, args, callback_func)
}

func (db *SQLMessagesDatabase) getMessagesWithCallback(ctx context.Context, where string, args []interface{}, callback_func GetMessagesCallbackFunc) error {

	pg_callback := func(pg_rsp pg_sql.PaginatedResponse) error {

		rows := pg_rsp.Rows()

		for rows.Next() {

			var id int64
			var note_id int64
			var author_address string
			var account_id int64
			var created int64
			var lastmod int64

			err := rows.Scan(&id, &note_id, &author_address, &account_id, &created, &lastmod)

			if err != nil {
				return fmt.Errorf("Failed to query database, %w", err)
			}

			m := &activitypub.Message{
				Id:            id,
				NoteId:        note_id,
				AuthorAddress: author_address,
				AccountId:     account_id,
				Created:       created,
				LastModified:  lastmod,
			}

			err = callback_func(ctx, m)

			if err != nil {
				return fmt.Errorf("Failed to execute following callback for message %d, %w", m.Id, err)
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

	q := fmt.Sprintf("SELECT id, note_id, author_address, account_id, created, lastmodified FROM %s WHERE %s ORDER BY created DESC", SQL_MESSAGES_TABLE_NAME, where)

	err = pg_sql.QueryPaginatedAll(db.database, pg_opts, pg_callback, q, args...)

	if err != nil {
		return fmt.Errorf("Failed to execute paginated query, %w", err)
	}

	return nil
}

func (db *SQLMessagesDatabase) Close(ctx context.Context) error {
	return db.database.Close()
}
