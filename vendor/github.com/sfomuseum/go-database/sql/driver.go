package sql

import (
	"database/sql"
	"fmt"
	"log/slog"
	"reflect"
)

const SQLITE_DRIVER string = "sqlite"
const MYSQL_DRIVER string = "mysql"
const POSTGRES_DRIVER string = "postgres"
const DUCKDB_DRIVER string = "duckdb"
const NULL_DRIVER string = "null"

// https://github.com/golang/go/issues/12600
// https://stackoverflow.com/questions/38811056/how-to-determine-name-of-database-driver-im-using

func DriverTypeOf(db *sql.DB) string {
	return fmt.Sprintf("%s", reflect.TypeOf(db.Driver()))
}

func Driver(db *sql.DB) string {

	driver_type := DriverTypeOf(db)

	switch driver_type {
	case "*sqlite3.SQLiteDriver", "*sqlite.Driver":
		return SQLITE_DRIVER
	case "*pq.Driver":
		return POSTGRES_DRIVER
	case "duckdb.Driver", "*duckdb.Driver":
		return DUCKDB_DRIVER
	case "*mysql.MySQLDriver", "mysql.MySQLDriver":
		return MYSQL_DRIVER
	case "*sql.nullDriver", "sql.nullDriver":
		return NULL_DRIVER
	default:
		slog.Warn("Unhandled driver type", "type", driver_type)
		return ""
	}
}
