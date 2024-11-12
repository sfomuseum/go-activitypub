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

const SQL_LIKES_TABLE_NAME string = "likes"

type SQLLikesDatabase struct {
	LikesDatabase
	database *sql.DB
}

func init() {

	ctx := context.Background()
	err := RegisterLikesDatabase(ctx, "sql", NewSQLLikesDatabase)

	if err != nil {
		panic(err)
	}
}

func NewSQLLikesDatabase(ctx context.Context, uri string) (LikesDatabase, error) {

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
			return nil, fmt.Errorf("Failed to configure SQLite, %w", err)
		}
	}

	db := &SQLLikesDatabase{
		database: conn,
	}

	return db, nil
}

func (db *SQLLikesDatabase) GetLikeIdsForDateRange(ctx context.Context, start int64, end int64, cb GetLikeIdsCallbackFunc) error {

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
				return fmt.Errorf("Failed to execute following callback for like %d, %w", id, err)
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

	q := fmt.Sprintf("SELECT id FROM %s WHERE created >= ? AND created <= ?", SQL_LIKES_TABLE_NAME)

	err = pg_sql.QueryPaginatedAll(db.database, pg_opts, pg_callback, q, start, end)

	if err != nil {
		return fmt.Errorf("Failed to execute paginated query, %w", err)
	}

	return nil
}

func (db *SQLLikesDatabase) GetLikeWithId(ctx context.Context, id int64) (*activitypub.Like, error) {

	where := "id = ?"
	return db.getLike(ctx, where, id)
}

func (db *SQLLikesDatabase) GetLikeWithPostIdAndActor(ctx context.Context, post_id int64, actor string) (*activitypub.Like, error) {

	where := "post_id = ? AND actor = ?"
	return db.getLike(ctx, where, post_id, actor)
}

func (db *SQLLikesDatabase) GetLikesForPostIdAndActor(ctx context.Context, post_id int64, actor string, cb GetLikesCallbackFunc) error {

	pg_callback := func(pg_rsp pg_sql.PaginatedResponse) error {

		rows := pg_rsp.Rows()

		for rows.Next() {

			var id int64
			var account_id int64
			var post_id int64
			var actor string
			var created int64

			err := rows.Scan(&id, &account_id, &post_id, &actor, &created)

			switch {
			case err == sql.ErrNoRows:
				return nil
			case err != nil:
				return err
			default:

				b := &activitypub.Like{
					Id:        id,
					AccountId: account_id,
					PostId:    post_id,
					Actor:     actor,
					Created:   created,
				}

				err = cb(ctx, b)

				if err != nil {
					return fmt.Errorf("Failed to execute callback for '%b', %w", b.Id, err)
				}

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

	q := fmt.Sprintf("SELECT id, account_id, post_id, actor, created FROM %s WHERE post_id = ? AND actor = ?", SQL_LIKES_TABLE_NAME)

	err = pg_sql.QueryPaginatedAll(db.database, pg_opts, pg_callback, q, post_id, actor)

	if err != nil {
		return fmt.Errorf("Failed to execute paginated query, %w", err)
	}

	return nil

}

func (db *SQLLikesDatabase) AddLike(ctx context.Context, b *activitypub.Like) error {

	q := fmt.Sprintf("INSERT INTO %s (id, account_id, psot_id, actor, created) VALUES (?, ?, ?, ?, ?)", SQL_LIKES_TABLE_NAME)

	_, err := db.database.ExecContext(ctx, q, b.Id, b.AccountId, b.PostId, b.Actor, b.Created)

	if err != nil {
		return fmt.Errorf("Failed to add like, %w", err)
	}

	return nil
}

func (db *SQLLikesDatabase) RemoveLike(ctx context.Context, b *activitypub.Like) error {

	q := fmt.Sprintf("DELETE FROM %s WHERE id= ?", SQL_LIKES_TABLE_NAME)

	_, err := db.database.ExecContext(ctx, q, b.Id)

	if err != nil {
		return fmt.Errorf("Failed to remove like, %w", err)
	}

	return nil
}

func (db *SQLLikesDatabase) Close(ctx context.Context) error {
	return db.database.Close()
}

func (db *SQLLikesDatabase) getLike(ctx context.Context, where string, args ...interface{}) (*activitypub.Like, error) {

	var id int64
	var account_id int64
	var post_id int64
	var actor string
	var created int64

	q := fmt.Sprintf("SELECT id, account_id, post_id, actor, created FROM %s WHERE %s", SQL_LIKES_TABLE_NAME, where)

	row := db.database.QueryRowContext(ctx, q, args...)

	err := row.Scan(&id, &account_id, &post_id, &actor, &created)

	switch {
	case err == sql.ErrNoRows:
		return nil, activitypub.ErrNotFound
	case err != nil:
		return nil, err
	default:
		//
	}

	b := &activitypub.Like{
		Id:        id,
		AccountId: account_id,
		PostId:    post_id,
		Actor:     actor,
		Created:   created,
	}

	return b, nil
}
