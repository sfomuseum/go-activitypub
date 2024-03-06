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

const SQL_DELIVERIES_TABLE_NAME string = "deliveries"

type SQLDeliveriesDatabase struct {
	DeliveriesDatabase
	database *sql.DB
}

func init() {
	ctx := context.Background()
	RegisterDeliveriesDatabase(ctx, "sql", NewSQLDeliveriesDatabase)
}

func NewSQLDeliveriesDatabase(ctx context.Context, uri string) (DeliveriesDatabase, error) {

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

	db := &SQLDeliveriesDatabase{
		database: conn,
	}

	return db, nil
}

func (db *SQLDeliveriesDatabase) GetDeliveryIdsForDateRange(ctx context.Context, start int64, end int64, cb GetDeliveryIdsCallbackFunc) error {

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
				return fmt.Errorf("Failed to execute following callback for delivery %d, %w", id, err)
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

	q := fmt.Sprintf("SELECT id FROM %s WHERE created >= ? AND created <= ?", SQL_DELIVERIES_TABLE_NAME)

	err = pg_sql.QueryPaginatedAll(db.database, pg_opts, pg_callback, q, start, end)

	if err != nil {
		return fmt.Errorf("Failed to execute paginated query, %w", err)
	}

	return nil
}

func (db *SQLDeliveriesDatabase) AddFollower(ctx context.Context, d *Delivery) error {

	/*
		q := fmt.Sprintf("INSERT INTO %s (id, account_id, follower_address, created) VALUES (?, ?, ?, ?)", SQL_DELIVERIES_TABLE_NAME)

		_, err := db.database.ExecContext(ctx, q, f.Id, f.AccountId, f.FollowerAddress, f.Created)

		if err != nil {
			return fmt.Errorf("Failed to add follower, %w", err)
		}
	*/

	return nil
}

func (db *SQLDeliveriesDatabase) Close(ctx context.Context) error {
	return db.database.Close()
}
