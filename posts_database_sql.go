package activitypub

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
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

	db := &SQLPostsDatabase{
		database: conn,
	}

	return db, nil
}

func (db *SQLPostsDatabase) AddPost(ctx context.Context, p *Post) error {

	q := fmt.Sprintf("INSERT INTO %s (id, account_id, body, created, lastmodified) VALUES (?, ?, ?, ?, ?)", SQL_POSTS_TABLE_NAME)

	_, err := db.database.ExecContext(ctx, q, p.Id, p.AccountId, p.Body, p.Created, p.LastModified)

	if err != nil {
		return fmt.Errorf("Failed to add post, %w", err)
	}

	return nil
}

func (db *SQLPostsDatabase) GetPost(ctx context.Context, id string) (*Post, error) {

	var account_id string
	var body []byte
	var created int64
	var lastmod int64

	q := fmt.Sprintf("SELECT account_id, body, created, lastmodified FROM %s WHERE id=?", SQL_POSTS_TABLE_NAME)

	row := db.database.QueryRowContext(ctx, q, id)

	err := row.Scan(&account_id, &body, &created, &lastmod)

	switch {
	case err == sql.ErrNoRows:
		return nil, err
	case err != nil:
		return nil, err
	default:
		//
	}

	a := &Post{
		Id:           id,
		AccountId:    account_id,
		Body:         body,
		Created:      created,
		LastModified: lastmod,
	}

	return a, nil
}