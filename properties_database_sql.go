package activitypub

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"

	"github.com/sfomuseum/go-activitypub/sqlite"
)

const SQL_PROPERTIES_TABLE_NAME string = "properties"

type SQLPropertiesDatabase struct {
	PropertiesDatabase
	database *sql.DB
}

func init() {
	ctx := context.Background()
	RegisterPropertiesDatabase(ctx, "sql", NewSQLPropertiesDatabase)
}

func NewSQLPropertiesDatabase(ctx context.Context, uri string) (PropertiesDatabase, error) {

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
			return nil, fmt.Errorf("Failed to set up SQLite, %w", err)
		}
	}

	db := &SQLPropertiesDatabase{
		database: conn,
	}

	return db, nil
}

func (db *SQLPropertiesDatabase) GetProperties(ctx context.Context, cb GetPropertiesCallbackFunc) error {
	return ErrNotImplemented
}

func (db *SQLPropertiesDatabase) GetPropertiesForAccount(ctx context.Context, account_id int64, cb GetPropertiesCallbackFunc) error {
	return ErrNotImplemented
}

func (db *SQLPropertiesDatabase) AddProperty(ctx context.Context, property *Property) error {
	return ErrNotImplemented
}

func (db *SQLPropertiesDatabase) UpdateProperty(ctx context.Context, property *Property) error {
	return ErrNotImplemented
}

func (db *SQLPropertiesDatabase) RemoveProperty(ctx context.Context, property *Property) error {
	return ErrNotImplemented
}

func (db *SQLPropertiesDatabase) Close(ctx context.Context) error {
	return db.database.Close()
}
