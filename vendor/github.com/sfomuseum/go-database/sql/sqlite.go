package sql

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
)

func DefaultSQLitePragma() []string {

	pragma := []string{
		"PRAGMA JOURNAL_MODE=OFF",
		"PRAGMA SYNCHRONOUS=OFF",
		// https://www.gaia-gis.it/gaia-sins/spatialite-cookbook/html/system.html
		"PRAGMA PAGE_SIZE=4096",
		"PRAGMA CACHE_SIZE=1000000",
	}

	return pragma
}

func ConfigureSQLitePragma(ctx context.Context, db *sql.DB, pragma []string) error {

	for _, p := range pragma {

		_, err := db.ExecContext(ctx, p)

		if err != nil {
			return fmt.Errorf("Failed to set pragma '%s', %w", p, err)
		}
	}

	return nil
}

func ConfigureSQLiteDatabase(ctx context.Context, db *sql.DB, opts *ConfigureDatabaseOptions) error {

	if opts.CreateTablesIfNecessary {

		table_names := make([]string, 0)

		sql := "SELECT name FROM sqlite_master WHERE type='table'"

		rows, err := db.QueryContext(ctx, sql)

		if err != nil {
			return fmt.Errorf("Failed to query sqlite_master, %w", err)
		}

		defer rows.Close()

		for rows.Next() {

			var name string
			err := rows.Scan(&name)

			if err != nil {
				return fmt.Errorf("Failed scan table name, %w", err)
			}

			table_names = append(table_names, name)
		}

		for _, t := range opts.Tables {

			if slices.Contains(table_names, t.Name()) {
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
		}
	}

	return nil
}

func HasSQLiteTable(ctx context.Context, db *sql.DB, table_name string) (bool, error) {

	has_table := false

	// TBD... how to derive database engine...

	sql := "SELECT name FROM sqlite_master WHERE type='table'"

	rows, err := db.QueryContext(ctx, sql)

	if err != nil {
		return false, fmt.Errorf("Failed to query sqlite_master, %w", err)
	}

	defer rows.Close()

	for rows.Next() {

		var name string
		err := rows.Scan(&name)

		if err != nil {
			return false, fmt.Errorf("Failed scan table name, %w", err)
		}

		if name == table_name {
			has_table = true
			break
		}
	}

	return has_table, nil
}
