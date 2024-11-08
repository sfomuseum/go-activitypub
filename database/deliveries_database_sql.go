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

const SQL_DELIVERIES_TABLE_NAME string = "deliveries"

type SQLDeliveriesDatabase struct {
	DeliveriesDatabase
	database *sql.DB
}

func init() {
	ctx := context.Background()
	err := RegisterDeliveriesDatabase(ctx, "sql", NewSQLDeliveriesDatabase)

	if err != nil {
		panic(err)
	}
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

func (db *SQLDeliveriesDatabase) AddDelivery(ctx context.Context, d *activitypub.Delivery) error {

	q := fmt.Sprintf("INSERT INTO %s (id, activity_id, activitypub_id, account_id, recipient, inbox, created, completed, success, error) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", SQL_DELIVERIES_TABLE_NAME)

	_, err := db.database.ExecContext(ctx, q, d.Id, d.ActivityId, d.ActivityPubId, d.AccountId, d.Recipient, d.Inbox, d.Created, d.Completed, d.Success, d.Error)

	if err != nil {
		return fmt.Errorf("Failed to add delivery, %w", err)
	}

	return nil
}

func (db *SQLDeliveriesDatabase) GetDeliveryWithId(ctx context.Context, id int64) (*activitypub.Delivery, error) {
	where := "id = ?"
	return db.getDelivery(ctx, where, id)
}

func (db *SQLDeliveriesDatabase) GetDeliveries(ctx context.Context, cb GetDeliveriesCallbackFunc) error {
	where := "1 = 1"
	args := make([]interface{}, 0)

	return db.getDeliveries(ctx, where, args, cb)
}

func (db *SQLDeliveriesDatabase) GetDeliveriesWithActivityIdAndRecipient(ctx context.Context, activity_id int64, recipient string, cb GetDeliveriesCallbackFunc) error {

	where := "activity_id = ? AND recipient = ?"
	args := []interface{}{
		activity_id,
		recipient,
	}

	return db.getDeliveries(ctx, where, args, cb)
}

func (db *SQLDeliveriesDatabase) GetDeliveriesWithActivityPubIdAndRecipient(ctx context.Context, activitypub_id string, recipient string, cb GetDeliveriesCallbackFunc) error {

	where := "activitypub_id = ? AND recipient = ?"
	args := []interface{}{
		activitypub_id,
		recipient,
	}

	return db.getDeliveries(ctx, where, args, cb)
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

func (db *SQLDeliveriesDatabase) Close(ctx context.Context) error {
	return db.database.Close()
}

func (db *SQLDeliveriesDatabase) getDelivery(ctx context.Context, where string, args ...interface{}) (*activitypub.Delivery, error) {

	var id int64
	var activity_id int64
	var activitypub_id string
	var account_id int64
	var recipient string
	var inbox string
	var created int64
	var completed int64
	var success_i int
	var error string

	q := fmt.Sprintf("SELECT id, activity_id, activitypub_id, account_id, recipient, inbox, created, completed, success, error FROM %s WHERE %s", SQL_DELIVERIES_TABLE_NAME, where)

	row := db.database.QueryRowContext(ctx, q, args...)

	err := row.Scan(&id, &activity_id, &activitypub_id, &account_id, &recipient, &inbox, &created, &completed, &success_i, &error)

	switch {
	case err == sql.ErrNoRows:
		return nil, activitypub.ErrNotFound
	case err != nil:
		return nil, err
	default:
		//
	}

	a := &activitypub.Delivery{
		Id:            id,
		ActivityId:    activity_id,
		ActivityPubId: activitypub_id,
		AccountId:     account_id,
		Recipient:     recipient,
		Inbox:         inbox,
		Created:       created,
		Completed:     completed,
		Error:         error,
	}

	if success_i > 0 {
		a.Success = true
	}

	return a, nil

}

func (db *SQLDeliveriesDatabase) getDeliveries(ctx context.Context, where string, args []interface{}, cb GetDeliveriesCallbackFunc) error {

	pg_callback := func(pg_rsp pg_sql.PaginatedResponse) error {

		rows := pg_rsp.Rows()

		for rows.Next() {

			var id int64
			var activity_id int64
			var activitypub_id string
			var account_id int64
			var recipient string
			var inbox string
			var created int64
			var completed int64
			var success_i int
			var error string

			err := rows.Scan(&id, &activity_id, &activitypub_id, &account_id, &recipient, &inbox, &created, &completed, &success_i, &error)

			if err != nil {
				return fmt.Errorf("Failed to query database, %w", err)
			}

			a := &activitypub.Delivery{
				Id:            id,
				ActivityId:    activity_id,
				ActivityPubId: activitypub_id,
				AccountId:     account_id,
				Recipient:     recipient,
				Inbox:         inbox,
				Created:       created,
				Completed:     completed,
				Error:         error,
			}

			if success_i > 0 {
				a.Success = true
			}

			err = cb(ctx, a)

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

	q := fmt.Sprintf("SELECT id, activity_id, activitypub_id, account_id, recipient, inbox, created, completed, success, error FROM %s WHERE %s", SQL_DELIVERIES_TABLE_NAME, where)

	err = pg_sql.QueryPaginatedAll(db.database, pg_opts, pg_callback, q, args...)

	if err != nil {
		return fmt.Errorf("Failed to execute paginated query, %w", err)
	}

	return nil

}
