package sql

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
)

var re_mem *regexp.Regexp
var re_vfs *regexp.Regexp
var re_file *regexp.Regexp

func init() {
	re_mem = regexp.MustCompile(`^(file\:)?\:memory\:.*`)
	re_vfs = regexp.MustCompile(`^vfs:\.*`)
	re_file = regexp.MustCompile(`^file\:([^\?]+)(?:\?.*)?$`)
}

type Table interface {
	Name() string
	Schema(*sql.DB) (string, error)
	InitializeTable(context.Context, *sql.DB) error
	IndexRecord(context.Context, *sql.DB, *sql.Tx, interface{}) error
}

func HasTable(ctx context.Context, db *sql.DB, table_name string) (bool, error) {

	switch Driver(db) {
	case SQLITE_DRIVER:
		return HasSQLiteTable(ctx, db, table_name)
	case POSTGRES_DRIVER:
		return HasPostgresTable(ctx, db, table_name)
	case DUCKDB_DRIVER:
		return HasDuckDBTable(ctx, db, table_name)
	case MYSQL_DRIVER:
		return HasMySQLTable(ctx, db, table_name)
	default:
		return false, fmt.Errorf("Unhandled or unsupported database driver %s", DriverTypeOf(db))
	}

}

func CreateTableIfNecessary(ctx context.Context, db *sql.DB, t Table) error {

	create := false

	has_table, err := HasTable(ctx, db, t.Name())

	if err != nil {
		return err
	}

	if !has_table {
		create = true
	}

	if create {

		sql, err := t.Schema(db)

		if err != nil {
			return err
		}

		_, err = db.ExecContext(ctx, sql)

		if err != nil {
			return err
		}

	}

	return nil
}

func IndexRecord(ctx context.Context, db *sql.DB, r interface{}, tables ...Table) error {

	tx, err := db.Begin()

	if err != nil {
		return err
	}

	for _, t := range tables {

		err := t.IndexRecord(ctx, db, tx, r)

		if err != nil {
			tx.Rollback()
			return fmt.Errorf("Failed to index %s table, %w", t.Name(), err)
		}
	}

	return tx.Commit()
}
