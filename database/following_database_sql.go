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

const SQL_FOLLOWING_TABLE_NAME string = "following"

type SQLFollowingDatabase struct {
	FollowingDatabase
	database *sql.DB
}

func init() {
	ctx := context.Background()
	err := RegisterFollowingDatabase(ctx, "sql", NewSQLFollowingDatabase)

	if err != nil {
		panic(err)
	}
}

func NewSQLFollowingDatabase(ctx context.Context, uri string) (FollowingDatabase, error) {

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

	db := &SQLFollowingDatabase{
		database: conn,
	}

	return db, nil
}

func (db *SQLFollowingDatabase) GetFollowingIdsForDateRange(ctx context.Context, start int64, end int64, cb GetFollowingIdsCallbackFunc) error {

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
				return fmt.Errorf("Failed to execute following callback for following %d, %w", id, err)
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

	q := fmt.Sprintf("SELECT id FROM %s WHERE created >= ? AND created <= ?", SQL_FOLLOWING_TABLE_NAME)

	err = pg_sql.QueryPaginatedAll(db.database, pg_opts, pg_callback, q, start, end)

	if err != nil {
		return fmt.Errorf("Failed to execute paginated query, %w", err)
	}

	return nil
}

func (db *SQLFollowingDatabase) GetFollowing(ctx context.Context, account_id int64, following_address string) (*activitypub.Following, error) {

	q := fmt.Sprintf("SELECT id, created FROM %s WHERE account_id = ? AND following_address = ?", SQL_FOLLOWING_TABLE_NAME)

	row := db.database.QueryRowContext(ctx, q, account_id, following_address)

	var id int64
	var created int64

	err := row.Scan(&id, &created)

	switch {
	case err == sql.ErrNoRows:
		return nil, activitypub.ErrNotFound
	case err != nil:
		return nil, fmt.Errorf("Failed to query database, %w", err)
	default:

		f := &activitypub.Following{
			Id:               id,
			AccountId:        account_id,
			FollowingAddress: following_address,
			Created:          created,
		}

		return f, nil
	}

}

func (db *SQLFollowingDatabase) AddFollowing(ctx context.Context, f *activitypub.Following) error {

	q := fmt.Sprintf("INSERT INTO %s (id, account_id, following_address, created) VALUES (?, ?, ?, ?)", SQL_FOLLOWING_TABLE_NAME)

	_, err := db.database.ExecContext(ctx, q, f.Id, f.AccountId, f.FollowingAddress, f.Created)

	if err != nil {
		return fmt.Errorf("Failed to add following, %w", err)
	}

	return nil
}

func (db *SQLFollowingDatabase) RemoveFollowing(ctx context.Context, f *activitypub.Following) error {

	q := fmt.Sprintf("DELETE FROM %s WHERE id = ?", SQL_FOLLOWING_TABLE_NAME)

	_, err := db.database.ExecContext(ctx, q, f.Id)

	if err != nil {
		return fmt.Errorf("Failed to remove following, %w", err)
	}

	return nil
}

func (db *SQLFollowingDatabase) GetFollowingForAccount(ctx context.Context, account_id int64, following_callback GetFollowingCallbackFunc) error {

	pg_callback := func(pg_rsp pg_sql.PaginatedResponse) error {

		rows := pg_rsp.Rows()

		for rows.Next() {

			var following_address string

			err := rows.Scan(&following_address)

			if err != nil {
				return fmt.Errorf("Failed to scan row, %w", err)
			}

			err = following_callback(ctx, following_address)

			if err != nil {
				return fmt.Errorf("Failed to execute following callback for '%s', %w", following_address, err)
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

	q := fmt.Sprintf("SELECT following_address FROM %s WHERE account_id=?", SQL_FOLLOWING_TABLE_NAME)

	err = pg_sql.QueryPaginatedAll(db.database, pg_opts, pg_callback, q, account_id)

	if err != nil {
		return fmt.Errorf("Failed to execute paginated query, %w", err)
	}

	return nil
}

func (db *SQLFollowingDatabase) Close(ctx context.Context) error {
	return db.database.Close()
}
