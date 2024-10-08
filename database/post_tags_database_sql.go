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

const SQL_POST_TAGS_TABLE_NAME string = "posts_tags"

type SQLPostTagsDatabase struct {
	PostTagsDatabase
	database *sql.DB
}

func init() {
	ctx := context.Background()
	RegisterPostTagsDatabase(ctx, "sql", NewSQLPostTagsDatabase)
}

func NewSQLPostTagsDatabase(ctx context.Context, uri string) (PostTagsDatabase, error) {

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

	db := &SQLPostTagsDatabase{
		database: conn,
	}

	return db, nil
}

func (db *SQLPostTagsDatabase) GetPostTagIdsForDateRange(ctx context.Context, start int64, end int64, cb GetPostTagIdsCallbackFunc) error {

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
				return fmt.Errorf("Failed to execute following callback for post tag %d, %w", id, err)
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

	q := fmt.Sprintf("SELECT id FROM %s WHERE created >= ? AND created <= ?", SQL_POST_TAGS_TABLE_NAME)

	err = pg_sql.QueryPaginatedAll(db.database, pg_opts, pg_callback, q, start, end)

	if err != nil {
		return fmt.Errorf("Failed to execute paginated query, %w", err)
	}

	return nil
}

func (db *SQLPostTagsDatabase) GetPostTagWithId(ctx context.Context, id int64) (*PostTag, error) {

	var account_id int64
	var post_id int64
	var href string
	var name string
	var pt_type string
	var created int64

	q := fmt.Sprintf("SELECT id, account_id, post_id, href, name, typem created FROM %s", SQL_POST_TAGS_TABLE_NAME)

	row := db.database.QueryRowContext(ctx, q, id)

	err := row.Scan(&id, &account_id, &post_id, &href, &name, &pt_type, &created)

	switch {
	case err == sql.ErrNoRows:
		return nil, ErrNotFound
	case err != nil:
		return nil, err
	default:

		t := &PostTag{
			Id:        id,
			AccountId: account_id,
			PostId:    post_id,
			Href:      href,
			Name:      name,
			Type:      pt_type,
			Created:   created,
		}

		return t, nil
	}
}

func (db *SQLPostTagsDatabase) GetPostTagsForName(ctx context.Context, name string, cb GetPostTagsCallbackFunc) error {

	where := "name = ?"

	args := []interface{}{
		name,
	}

	return db.getPostTagsWithCallback(ctx, where, args, cb)
}

func (db *SQLPostTagsDatabase) GetPostTagsForAccount(ctx context.Context, account_id int64, cb GetPostTagsCallbackFunc) error {

	where := "account_id = ?"

	args := []interface{}{
		account_id,
	}

	return db.getPostTagsWithCallback(ctx, where, args, cb)
}

func (db *SQLPostTagsDatabase) GetPostTagsForPost(ctx context.Context, post_id int64, cb GetPostTagsCallbackFunc) error {

	where := "post_id = ?"

	args := []interface{}{
		post_id,
	}

	return db.getPostTagsWithCallback(ctx, where, args, cb)
}

func (db *SQLPostTagsDatabase) AddPostTag(ctx context.Context, post_tag *PostTag) error {

	q := fmt.Sprintf("INSERT INTO %s (id, account_id, post_id, href, name, type, created) VALUES (?, ?, ?, ?, ?, ?, ?)", SQL_POST_TAGS_TABLE_NAME)

	_, err := db.database.ExecContext(ctx, q, post_tag.Id, post_tag.AccountId, post_tag.PostId, post_tag.Href, post_tag.Name, post_tag.Type, post_tag.Created)

	if err != nil {
		return fmt.Errorf("Failed to add post tag, %w", err)
	}

	return nil
}

func (db *SQLPostTagsDatabase) RemovePostTag(ctx context.Context, post_tag *PostTag) error {

	q := fmt.Sprintf("DELETE FROM %s WHERE id = ?", SQL_POST_TAGS_TABLE_NAME)

	_, err := db.database.ExecContext(ctx, q, post_tag.Id)

	if err != nil {
		return fmt.Errorf("Failed to remove post tag, %w", err)
	}

	return nil
}

func (db *SQLPostTagsDatabase) Close(ctx context.Context) error {
	return db.database.Close()
}

func (db *SQLPostTagsDatabase) getPostTagsWithCallback(ctx context.Context, where string, args []interface{}, callback_func GetPostTagsCallbackFunc) error {

	pg_callback := func(pg_rsp pg_sql.PaginatedResponse) error {

		rows := pg_rsp.Rows()

		for rows.Next() {

			var id int64
			var account_id int64
			var post_id int64
			var href string
			var name string
			var pt_type string
			var created int64

			err := rows.Scan(&id, &account_id, &post_id, &href, &name, &pt_type, &created)

			if err != nil {
				return fmt.Errorf("Failed to query database, %w", err)
			}

			t := &PostTag{
				Id:        id,
				AccountId: account_id,
				PostId:    post_id,
				Href:      href,
				Name:      name,
				Type:      pt_type,
				Created:   created,
			}

			err = callback_func(ctx, t)

			if err != nil {
				return fmt.Errorf("Failed to execute following callback for post tag %d, %w", t.Id, err)
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

	q := fmt.Sprintf("SELECT id, account_id, post_id, href, name, typem created FROM %s WHERE %s ORDER BY created DESC", SQL_POST_TAGS_TABLE_NAME, where)

	err = pg_sql.QueryPaginatedAll(db.database, pg_opts, pg_callback, q, args...)

	if err != nil {
		return fmt.Errorf("Failed to execute paginated query, %w", err)
	}

	return nil
}
