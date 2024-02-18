package activitypub

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"

	"github.com/sfomuseum/go-activitypub/sqlite"
)

const SQL_MESSAGES_TABLE_NAME string = "messages"

type SQLMessagesDatabase struct {
	MessagesDatabase
	database *sql.DB
}

func init() {
	ctx := context.Background()
	RegisterMessagesDatabase(ctx, "sql", NewSQLMessagesDatabase)
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

func (db *SQLMessagesDatabase) GetMessageWithId(ctx context.Context, message_id int64) (*Message, error) {
	where := "id = ?"
	return db.getMessage(ctx, where, message_id)
}

func (db *SQLMessagesDatabase) GetMessageWithAccountAndNoteIds(ctx context.Context, account_id int64, note_id int64) (*Message, error) {

	where := "account_id = ? AND note_id = ?"
	return db.getMessage(ctx, where, account_id, note_id)

}

func (db *SQLMessagesDatabase) getMessage(ctx context.Context, where string, args ...interface{}) (*Message, error) {

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
		return nil, ErrNotFound
	case err != nil:
		return nil, fmt.Errorf("Failed to query database, %w", err)
	default:

		n := &Message{
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

func (db *SQLMessagesDatabase) AddMessage(ctx context.Context, message *Message) error {

	q := fmt.Sprintf("INSERT INTO %s (id, note_id, author_address, account_id, created, lastmodified) VALUES (?, ?, ?, ?, ?, ?)", SQL_MESSAGES_TABLE_NAME)

	_, err := db.database.ExecContext(ctx, q, message.Id, message.NoteId, message.AuthorAddress, message.AccountId, message.Created, message.LastModified)

	if err != nil {
		return fmt.Errorf("Failed to add message, %w", err)
	}

	return nil
}

func (db *SQLMessagesDatabase) UpdateMessage(ctx context.Context, message *Message) error {

	q := fmt.Sprintf("UPDATE %s SET note_id=?, author_address=?, account_id=?, created=?, lastmodified=? WHERE id = ?", SQL_MESSAGES_TABLE_NAME)

	_, err := db.database.ExecContext(ctx, q, message.NoteId, message.AuthorAddress, message.AccountId, message.Created, message.LastModified, message.Id)

	if err != nil {
		return fmt.Errorf("Failed to add message, %w", err)
	}

	return nil
}

func (db *SQLMessagesDatabase) RemoveMessage(ctx context.Context, message *Message) error {

	q := fmt.Sprintf("DELETE FROM %s WHERE id = ?", SQL_MESSAGES_TABLE_NAME)

	_, err := db.database.ExecContext(ctx, q, message.Id)

	if err != nil {
		return fmt.Errorf("Failed to remove message, %w", err)
	}

	return nil
}
