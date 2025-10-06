package sql

import (
	"context"
	"database/sql"
	"fmt"
	_ "log/slog"
	"net/url"
)

type ConfigureDatabaseOptions struct {
	CreateTablesIfNecessary bool
	Tables                  []Table
	Pragma                  []string
}

func DefaultConfigureDatabaseOptions() *ConfigureDatabaseOptions {
	opts := &ConfigureDatabaseOptions{}
	return opts
}

func ConfigureDatabase(ctx context.Context, db *sql.DB, opts *ConfigureDatabaseOptions) error {

	switch Driver(db) {
	case SQLITE_DRIVER:
		return ConfigureSQLiteDatabase(ctx, db, opts)
	case POSTGRES_DRIVER:
		return ConfigurePostgresDatabase(ctx, db, opts)
	case DUCKDB_DRIVER:
		return ConfigureDuckDBDatabase(ctx, db, opts)
	default:
		return fmt.Errorf("Unhandled or unsupported database driver %s", DriverTypeOf(db))
	}
}

func OpenWithURI(ctx context.Context, db_uri string) (*sql.DB, error) {

	u, err := url.Parse(db_uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	q := u.Query()

	engine := u.Host
	dsn := q.Get("dsn")

	if engine == "" {
		return nil, fmt.Errorf("Missing database engine")
	}

	if dsn == "" {
		return nil, fmt.Errorf("Missing DSN string")
	}

	db, err := sql.Open(engine, dsn)

	if err != nil {
		return nil, fmt.Errorf("Unable to create database (%s) because %v", engine, err)
	}

	switch Driver(db) {
	case "sqlite":

		pragma := DefaultSQLitePragma()
		err := ConfigureSQLitePragma(ctx, db, pragma)

		if err != nil {
			return nil, fmt.Errorf("Failed to configure SQLite pragma, %w", err)
		}
	}

	return db, nil
}
