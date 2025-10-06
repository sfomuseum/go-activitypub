package database

import (
	"context"
	"database/sql"
	"fmt"

	pg_sql "github.com/aaronland/go-pagination-sql"
	"github.com/aaronland/go-pagination/countable"
	"github.com/sfomuseum/go-activitypub"
	sfom_sql "github.com/sfomuseum/go-database/sql"
)

const SQL_PROPERTIES_TABLE_NAME string = "properties"

type SQLPropertiesDatabase struct {
	PropertiesDatabase
	database *sql.DB
}

func init() {
	ctx := context.Background()
	err := RegisterPropertiesDatabase(ctx, "sql", NewSQLPropertiesDatabase)

	if err != nil {
		panic(err)
	}
}

func NewSQLPropertiesDatabase(ctx context.Context, uri string) (PropertiesDatabase, error) {

	conn, err := sfom_sql.OpenWithURI(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to open database connection, %w", err)
	}

	db := &SQLPropertiesDatabase{
		database: conn,
	}

	return db, nil
}

func (db *SQLPropertiesDatabase) GetProperties(ctx context.Context, cb GetPropertiesCallbackFunc) error {

	q := fmt.Sprintf("SELECT id, account_id, key, value, created FROM %s", SQL_PROPERTIES_TABLE_NAME)
	return db.getProperties(ctx, cb, q)
}

func (db *SQLPropertiesDatabase) GetPropertiesForAccount(ctx context.Context, account_id int64, cb GetPropertiesCallbackFunc) error {

	q := fmt.Sprintf("SELECT id, account_id, key, value, created FROM %s WHERE account_id=?", SQL_PROPERTIES_TABLE_NAME)
	return db.getProperties(ctx, cb, q, account_id)
}

func (db *SQLPropertiesDatabase) AddProperty(ctx context.Context, p *activitypub.Property) error {

	q := fmt.Sprintf("INSERT INTO %s (id, account_id, key, value, created) VALUES (?, ?, ?, ?, ?)", SQL_PROPERTIES_TABLE_NAME)

	_, err := db.database.ExecContext(ctx, q, p.Id, p.AccountId, p.Key, p.Value, p.Created)

	if err != nil {
		return fmt.Errorf("Failed to add property, %w", err)
	}

	return nil
}

func (db *SQLPropertiesDatabase) UpdateProperty(ctx context.Context, p *activitypub.Property) error {

	q := fmt.Sprintf("UPDATE %s SET account_id=?, key=?, value=?, created=? WHERE id=?", SQL_PROPERTIES_TABLE_NAME)

	_, err := db.database.ExecContext(ctx, q, p.AccountId, p.Key, p.Value, p.Created, p.Id)

	if err != nil {
		return fmt.Errorf("Failed to add property, %w", err)
	}

	return nil
}

func (db *SQLPropertiesDatabase) RemoveProperty(ctx context.Context, p *activitypub.Property) error {

	q := fmt.Sprintf("DELETE FROM %s WHERE id = ?", SQL_PROPERTIES_TABLE_NAME)

	_, err := db.database.ExecContext(ctx, q, p.Id)

	if err != nil {
		return fmt.Errorf("Failed to add property, %w", err)
	}

	return nil
}

func (db *SQLPropertiesDatabase) Close(ctx context.Context) error {
	return db.database.Close()
}

func (db *SQLPropertiesDatabase) getProperties(ctx context.Context, cb GetPropertiesCallbackFunc, q string, args ...any) error {

	pg_callback := func(pg_rsp pg_sql.PaginatedResponse) error {

		rows := pg_rsp.Rows()

		for rows.Next() {

			var id int64
			var account_id int64
			var key string
			var value string
			var created int64

			err := rows.Scan(&id, &account_id, &key, &value, &created)

			if err != nil {
				return fmt.Errorf("Failed to scan database row, %w", err)
			}

			pr := &activitypub.Property{
				Id:        id,
				AccountId: account_id,
				Key:       key,
				Value:     value,
				Created:   created,
			}

			err = cb(ctx, pr)

			if err != nil {
				return fmt.Errorf("Failed to execute callback for property %d, %w", id, err)
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

	err = pg_sql.QueryPaginatedAll(db.database, pg_opts, pg_callback, q, args...)

	if err != nil {
		return fmt.Errorf("Failed to execute paginated query, %w", err)
	}

	return nil
}
