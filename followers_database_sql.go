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

const SQL_FOLLOWERS_TABLE_NAME string = "followers"

type SQLFollowersDatabase struct {
	FollowersDatabase
	database *sql.DB
}

func init() {
	ctx := context.Background()
	RegisterFollowersDatabase(ctx, "sql", NewSQLFollowersDatabase)
}

func NewSQLFollowersDatabase(ctx context.Context, uri string) (FollowersDatabase, error) {

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

		conn.SetMaxOpenConns(1)

		err := sqlite.SetupConnection(ctx, conn)

		if err != nil {
			return nil, fmt.Errorf("Failed to configure SQLite, %w", err)
		}
	}

	db := &SQLFollowersDatabase{
		database: conn,
	}

	return db, nil
}

func (db *SQLFollowersDatabase) GetFollowerIdsForDateRange(ctx context.Context, start int64, end int64, cb GetFollowerIdsCallbackFunc) error {

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
				return fmt.Errorf("Failed to execute following callback for follower %d, %w", id, err)
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

	q := fmt.Sprintf("SELECT id FROM %s WHERE created >= ? AND created <= ?", SQL_FOLLOWERS_TABLE_NAME)

	err = pg_sql.QueryPaginatedAll(db.database, pg_opts, pg_callback, q, start, end)

	if err != nil {
		return fmt.Errorf("Failed to execute paginated query, %w", err)
	}

	return nil
}

func (db *SQLFollowersDatabase) HasFollower(ctx context.Context, account_id int64, follower_address string) (bool, error) {

	q := fmt.Sprintf("SELECT COUNT(id) AS count FROM %s WHERE account_id = ?", SQL_FOLLOWERS_TABLE_NAME)

	row := db.database.QueryRowContext(ctx, q, account_id)

	var count int

	err := row.Scan(&count)

	switch {
	case err == sql.ErrNoRows:
		return false, nil
	case err != nil:
		return false, fmt.Errorf("Failed to query database, %w", err)
	default:
		if count > 0 {
			return true, nil
		}
	}

	return false, nil
}

func (db *SQLFollowersDatabase) GetFollower(ctx context.Context, account_id int64, follower_address string) (*Follower, error) {

	q := fmt.Sprintf("SELECT id, created FROM %s WHERE account_id = ? AND follower_address = ?", SQL_FOLLOWERS_TABLE_NAME)

	row := db.database.QueryRowContext(ctx, q, account_id, follower_address)

	var id int64
	var created int64

	err := row.Scan(&id, &created)

	switch {
	case err == sql.ErrNoRows:
		return nil, ErrNotFound
	case err != nil:
		return nil, fmt.Errorf("Failed to query database, %w", err)
	default:

		f := &Follower{
			Id:              id,
			AccountId:       account_id,
			FollowerAddress: follower_address,
			Created:         created,
		}

		return f, nil
	}

}

func (db *SQLFollowersDatabase) AddFollower(ctx context.Context, f *Follower) error {

	q := fmt.Sprintf("INSERT INTO %s (id, account_id, follower_address, created) VALUES (?, ?, ?, ?)", SQL_FOLLOWERS_TABLE_NAME)

	_, err := db.database.ExecContext(ctx, q, f.Id, f.AccountId, f.FollowerAddress, f.Created)

	if err != nil {
		return fmt.Errorf("Failed to add follower, %w", err)
	}

	return nil
}

func (db *SQLFollowersDatabase) RemoveFollower(ctx context.Context, f *Follower) error {

	q := fmt.Sprintf("DELETE FROM %s WHERE id = ?", SQL_FOLLOWERS_TABLE_NAME)

	_, err := db.database.ExecContext(ctx, q, f.Id)

	if err != nil {
		return fmt.Errorf("Failed to remove follower, %w", err)
	}

	return nil
}

func (db *SQLFollowersDatabase) GetFollowersForAccount(ctx context.Context, account_id int64, followers_callback GetFollowersCallbackFunc) error {

	pg_callback := func(pg_rsp pg_sql.PaginatedResponse) error {

		rows := pg_rsp.Rows()

		for rows.Next() {

			var follower_address string

			err := rows.Scan(&follower_address)

			if err != nil {
				return fmt.Errorf("Failed to scan row, %w", err)
			}

			err = followers_callback(ctx, follower_address)

			if err != nil {
				return fmt.Errorf("Failed to execute followers callback for '%s', %w", follower_address, err)
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

	q := fmt.Sprintf("SELECT follower_address FROM %s WHERE account_id=?", SQL_FOLLOWERS_TABLE_NAME)

	err = pg_sql.QueryPaginatedAll(db.database, pg_opts, pg_callback, q, account_id)

	if err != nil {
		return fmt.Errorf("Failed to execute paginated query, %w", err)
	}

	return nil
}

func (db *SQLFollowersDatabase) Close(ctx context.Context) error {
	return db.database.Close()
}
