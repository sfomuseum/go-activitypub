package activitypub

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"

	pg_sql "github.com/aaronland/go-pagination-sql"
	"github.com/aaronland/go-pagination/countable"
	"github.com/sfomuseum/go-activitypub/sqlite"
)

const SQL_ALIASES_TABLE_NAME string = "aliases"

type SQLAliasesDatabase struct {
	AliasesDatabase
	database *sql.DB
}

func init() {
	ctx := context.Background()
	RegisterAliasesDatabase(ctx, "sql", NewSQLAliasesDatabase)
}

func NewSQLAliasesDatabase(ctx context.Context, uri string) (AliasesDatabase, error) {

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

	db := &SQLAliasesDatabase{
		database: conn,
	}

	return db, nil
}

func (db *SQLAliasesDatabase) GetAliasesForAccount(ctx context.Context, account_id int64, cb GetAliasesCallbackFunc) error {

	where := "account_id = ?"

	args := []interface{}{
		account_id,
	}

	return db.getAliasesWithCallback(ctx, where, args, cb)
}

func (db *SQLAliasesDatabase) GetAliasWithName(ctx context.Context, name string) (*Alias, error) {

	where := "name = ?"

	args := []interface{}{
		name,
	}

	return db.getAlias(ctx, where, args)
}

func (db *SQLAliasesDatabase) AddAlias(ctx context.Context, alias *Alias) error {

	q := fmt.Sprintf("INSERT INTO %s (name, account_id, created) VALUES (?, ?, ?)", SQL_ALIASES_TABLE_NAME)

	_, err := db.database.ExecContext(ctx, q, alias.Name, alias.AccountId, alias.Created)

	if err != nil {
		return fmt.Errorf("Failed to add post tag, %w", err)
	}

	return nil
}

func (db *SQLAliasesDatabase) RemoveAlias(ctx context.Context, alias *Alias) error {

	q := fmt.Sprintf("DELETE FROM %s WHERE name = ?", SQL_ALIASES_TABLE_NAME)

	_, err := db.database.ExecContext(ctx, q, alias.Name)

	if err != nil {
		return fmt.Errorf("Failed to remove post tag, %w", err)
	}

	return nil
}

func (db *SQLAliasesDatabase) Close(ctx context.Context) error {
	return db.database.Close()
}

func (db *SQLAliasesDatabase) getAlias(ctx context.Context, where string, args []interface{}) (*Alias, error) {

	var name string
	var account_id int64
	var created int64

	q := fmt.Sprintf("SELECT name, account_id, created FROM %s WHERE %s", SQL_ALIASES_TABLE_NAME, where)

	row := db.database.QueryRowContext(ctx, q, args...)

	err := row.Scan(&name, &account_id, &created)

	switch {
	case err == sql.ErrNoRows:
		return nil, ErrNotFound
	case err != nil:
		return nil, err
	default:

		a := &Alias{
			Name:      name,
			AccountId: account_id,
			Created:   created,
		}

		return a, nil
	}
}

func (db *SQLAliasesDatabase) getAliasesWithCallback(ctx context.Context, where string, args []interface{}, callback_func GetAliasesCallbackFunc) error {

	pg_callback := func(pg_rsp pg_sql.PaginatedResponse) error {

		rows := pg_rsp.Rows()

		for rows.Next() {

			var name string
			var account_id int64
			var created int64

			err := rows.Scan(&name, &account_id, &created)

			if err != nil {
				return fmt.Errorf("Failed to query database, %w", err)
			}

			a := &Alias{
				Name:      name,
				AccountId: account_id,
				Created:   created,
			}

			err = callback_func(ctx, a)

			if err != nil {
				return fmt.Errorf("Failed to execute following callback for alias %s, %w", a.Name, err)
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

	q := fmt.Sprintf("SELECT name, account_id, created FROM %s WHERE %s ORDER BY created DESC", SQL_ALIASES_TABLE_NAME, where)

	err = pg_sql.QueryPaginatedAll(db.database, pg_opts, pg_callback, q, args...)

	if err != nil {
		return fmt.Errorf("Failed to execute paginated query, %w", err)
	}

	return nil
}
