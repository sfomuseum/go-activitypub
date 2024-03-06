package activitypub

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"

	pg_sql "github.com/aaronland/go-pagination-sql"
	"github.com/aaronland/go-pagination/countable"
	"github.com/sfomuseum/go-activitypub/sqlite"
)

const SQL_POSTS_TABLE_NAME string = "posts"

type SQLPostsDatabase struct {
	PostsDatabase
	database *sql.DB
}

func init() {
	ctx := context.Background()
	RegisterPostsDatabase(ctx, "sql", NewSQLPostsDatabase)
}

func NewSQLPostsDatabase(ctx context.Context, uri string) (PostsDatabase, error) {

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

	db := &SQLPostsDatabase{
		database: conn,
	}

	return db, nil
}

func (db *SQLPostsDatabase) GetPostIdsForDateRange(ctx context.Context, start int64, end int64, cb GetPostIdsCallbackFunc) error {

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
				return fmt.Errorf("Failed to execute following callback for post %d, %w", id, err)
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

	q := fmt.Sprintf("SELECT id FROM %s WHERE created >= ? AND created <= ?", SQL_POSTS_TABLE_NAME)

	err = pg_sql.QueryPaginatedAll(db.database, pg_opts, pg_callback, q, start, end)

	if err != nil {
		return fmt.Errorf("Failed to execute paginated query, %w", err)
	}

	return nil
}

func (db *SQLPostsDatabase) AddPost(ctx context.Context, p *Post) error {

	q := fmt.Sprintf("INSERT INTO %s (id, account_id, body, in_reply_to, created, lastmodified) VALUES (?, ?, ?, ?, ?, ?)", SQL_POSTS_TABLE_NAME)

	_, err := db.database.ExecContext(ctx, q, p.Id, p.AccountId, p.Body, p.InReplyTo, p.Created, p.LastModified)

	if err != nil {
		return fmt.Errorf("Failed to add post, %w", err)
	}

	return nil
}

func (db *SQLPostsDatabase) GetPostWithId(ctx context.Context, id int64) (*Post, error) {
	where := "id = ?"
	return db.getPost(ctx, where, id)
}

func (db *SQLPostsDatabase) getPost(ctx context.Context, where string, args ...interface{}) (*Post, error) {

	var id int64
	var account_id int64
	var body string
	var in_reply_to string
	var created int64
	var lastmod int64

	q := fmt.Sprintf("SELECT id, account_id, body, in_reply_to, created, lastmodified FROM %s WHERE %s", SQL_POSTS_TABLE_NAME, where)

	row := db.database.QueryRowContext(ctx, q, id)

	err := row.Scan(&id, &account_id, &body, &in_reply_to, &created, &lastmod)

	switch {
	case err == sql.ErrNoRows:
		return nil, ErrNotFound
	case err != nil:
		return nil, err
	default:
		//
	}

	a := &Post{
		Id:           id,
		AccountId:    account_id,
		Body:         body,
		InReplyTo:    in_reply_to,
		Created:      created,
		LastModified: lastmod,
	}

	return a, nil
}

func (db *SQLPostsDatabase) Close(ctx context.Context) error {
	return db.database.Close()
}
