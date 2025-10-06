package sql

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
)

func ConfigurePostgresDatabase(ctx context.Context, db *sql.DB, opts *ConfigureDatabaseOptions) error {

	if opts.CreateTablesIfNecessary {

		for _, t := range opts.Tables {

			logger := slog.Default()
			logger = logger.With("table", t.Name())

			exists, err := HasPostgresTable(ctx, db, t.Name())

			if err != nil {
				return fmt.Errorf("Failed to determine if table %s exists, %w", t.Name(), err)
			}

			if exists {
				logger.Debug("Table already exists")
				continue
			}

			schema, err := t.Schema(db)

			if err != nil {
				return fmt.Errorf("Failed to derive schema for table %s, %w", t.Name(), err)
			}

			_, err = db.ExecContext(ctx, schema)

			if err != nil {
				return fmt.Errorf("Failed to create %s table, %w", t.Name(), err)
			}

			logger.Debug("Created table")
		}
	}

	return nil
}

// https://stackoverflow.com/questions/20582500/how-to-check-if-a-table-exists-in-a-given-schema

func HasPostgresTable(ctx context.Context, db *sql.DB, table_name string) (bool, error) {

	logger := slog.Default()
	logger = logger.With("table", table_name)

	q := "SELECT EXISTS(SELECT * FROM pg_tables WHERE schemaname='public' AND tablename=$1)"

	row := db.QueryRowContext(ctx, q, table_name)

	var exists bool

	err := row.Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("Failed to query table, %w", err)
	}

	logger.Debug("Does table exist", "exists", exists)
	return exists, nil
}
