package sql

import (
	"context"
	db_sql "database/sql"
	"fmt"
	"log/slog"
)

func ConfigureMySQLDatabase(ctx context.Context, db *db_sql.DB, opts *ConfigureDatabaseOptions) error {
	return nil
}

func HasMySQLTable(ctx context.Context, db *db_sql.DB, table_name string) (bool, error) {

	logger := slog.Default()
	logger = logger.With("table", table_name)

	exists := false

	q := "SHOW TABLES"

	rows, err := db.QueryContext(ctx, q)

	if err != nil {
		return false, fmt.Errorf("Failed to show tables, %w", err)
	}

	defer rows.Close()

	for rows.Next() {

		var db_table string
		err := rows.Scan(&db_table)

		if err != nil {
			return false, fmt.Errorf("Failed to scan row, %w", err)
		}

		if db_table == table_name {
			exists = true
			break
		}
	}

	logger.Debug("Does table exist", "exists", exists)
	return exists, nil
}
