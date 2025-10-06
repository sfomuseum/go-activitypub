package database

import (
	"context"
	"database/sql"
	"fmt"

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
	return activitypub.ErrNotImplemented
}

func (db *SQLPropertiesDatabase) GetPropertiesForAccount(ctx context.Context, account_id int64, cb GetPropertiesCallbackFunc) error {
	return activitypub.ErrNotImplemented
}

func (db *SQLPropertiesDatabase) AddProperty(ctx context.Context, property *activitypub.Property) error {
	return activitypub.ErrNotImplemented
}

func (db *SQLPropertiesDatabase) UpdateProperty(ctx context.Context, property *activitypub.Property) error {
	return activitypub.ErrNotImplemented
}

func (db *SQLPropertiesDatabase) RemoveProperty(ctx context.Context, property *activitypub.Property) error {
	return activitypub.ErrNotImplemented
}

func (db *SQLPropertiesDatabase) Close(ctx context.Context) error {
	return db.database.Close()
}
