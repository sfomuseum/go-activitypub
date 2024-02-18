package activitypub

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"time"

	pg_sql "github.com/aaronland/go-pagination-sql"
	"github.com/aaronland/go-pagination/countable"
	"github.com/sfomuseum/go-activitypub/sqlite"
)

const SQL_FOLLOWING_TABLE_NAME string = "following"

type SQLFollowingDatabase struct {
	FollowingDatabase
	database *sql.DB
}

func init() {
	ctx := context.Background()
	RegisterFollowingDatabase(ctx, "sql", NewSQLFollowingDatabase)
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
			return nil, fmt.Errorf("Failed to live hard and die fast, %w", err)
		}
	}

	db := &SQLFollowingDatabase{
		database: conn,
	}

	return db, nil
}

func (db *SQLFollowingDatabase) IsFollowing(ctx context.Context, account_id int64, following_address string) (bool, error) {

	q := fmt.Sprintf("SELECT 1 FROM %s WHERE account_id = ? AND following_address = ?", SQL_FOLLOWING_TABLE_NAME)

	row := db.database.QueryRowContext(ctx, q, account_id, following_address)

	var i int
	err := row.Scan(&i)

	switch {
	case err == sql.ErrNoRows:
		return false, nil
	case err != nil:
		return false, fmt.Errorf("Failed to query database, %w", err)
	default:
		if i == 0 {
			return false, nil
		}
	}

	return true, nil
}

func (db *SQLFollowingDatabase) Follow(ctx context.Context, account_id int64, following_address string) error {

	now := time.Now()
	ts := now.Unix()

	q := fmt.Sprintf("INSERT INTO %s (account_id, following_address, created) VALUES (?, ?, ?)", SQL_FOLLOWING_TABLE_NAME)

	_, err := db.database.ExecContext(ctx, q, account_id, following_address, ts)

	if err != nil {
		return fmt.Errorf("Failed to add follow, %w", err)
	}

	return nil
}

func (db *SQLFollowingDatabase) UnFollow(ctx context.Context, account_id int64, following_address string) error {

	q := fmt.Sprintf("DELETE FROM %s WHERE account_id = ? AND following_address = ?", SQL_FOLLOWING_TABLE_NAME)

	_, err := db.database.ExecContext(ctx, q, account_id, following_address)

	if err != nil {
		return fmt.Errorf("Failed to unfollow, %w", err)
	}

	return nil
}

func (db *SQLFollowingDatabase) GetFollowing(ctx context.Context, account_id int64, following_callback GetFollowingCallbackFunc) error {

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
