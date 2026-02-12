package database

import (
	"context"
	"database/sql"
	"fmt"
	"iter"

	pg_sql "github.com/aaronland/go-pagination-sql"
	"github.com/aaronland/go-pagination/countable"
	"github.com/sfomuseum/go-activitypub"
	sfom_sql "github.com/sfomuseum/go-database/sql"
)

const SQL_ACCOUNTS_TABLE_NAME string = "accounts"

type SQLAccountsDatabase struct {
	Database[*activitypub.Account]
	AccountsDatabase
	database *sql.DB
}

func init() {
	ctx := context.Background()
	err := RegisterAccountsDatabase(ctx, "sql", NewSQLAccountsDatabase)

	if err != nil {
		panic(err)
	}
}

func NewSQLAccountsDatabase(ctx context.Context, uri string) (AccountsDatabase, error) {

	conn, err := sfom_sql.OpenWithURI(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to open database connection, %w", err)
	}

	db := &SQLAccountsDatabase{
		database: conn,
	}

	return db, nil
}

func (db *SQLAccountsDatabase) GetRecord(ctx context.Context, id int64) (*activitypub.Account, error) {
	where := "id = ?"
	return db.getAccount(ctx, where, id)
}

func (db *SQLAccountsDatabase) AddRecord(ctx context.Context, a *activitypub.Account) error {

	q := fmt.Sprintf("INSERT INTO %s (id, account_type, name, display_name, blurb, url, public_key_uri, private_key_uri, created, lastmodified) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", SQL_ACCOUNTS_TABLE_NAME)

	_, err := db.database.ExecContext(ctx, q, a.Id, a.AccountType, a.Name, a.DisplayName, a.Blurb, a.URL, a.PublicKeyURI, a.PrivateKeyURI, a.Created, a.LastModified)

	if err != nil {
		return fmt.Errorf("Failed to add account, %w", err)
	}

	return nil
}

func (db *SQLAccountsDatabase) UpdateRecord(ctx context.Context, a *activitypub.Account) error {

	q := fmt.Sprintf("UPDATE %s SET account_type = ?, name = ?, display_name = ?, blurb = ?, url = ?, public_key_uri = ?, private_key_uri = ?, lastmodified =? WHERE id = ?", SQL_ACCOUNTS_TABLE_NAME)

	_, err := db.database.ExecContext(ctx, q, a.AccountType, a.Name, a.DisplayName, a.Blurb, a.URL, a.PublicKeyURI, a.PrivateKeyURI, a.LastModified, a.Id)

	if err != nil {
		return fmt.Errorf("Failed to update account, %w", err)
	}

	return nil
}

func (db *SQLAccountsDatabase) RemoveRecord(ctx context.Context, a *activitypub.Account) error {

	q := fmt.Sprintf("DELETE FROM %s WHERE id = ?", SQL_ACCOUNTS_TABLE_NAME)

	_, err := db.database.ExecContext(ctx, q, a.Id)

	if err != nil {
		return fmt.Errorf("Failed to delete account, %w", err)
	}

	return nil
}

func (db *SQLAccountsDatabase) QueryRecords(ctx context.Context, q *Query) iter.Seq2[*activitypub.Account, error] {

	return func(yield func(*activitypub.Account, error) bool) {

		pg_callback := func(pg_rsp pg_sql.PaginatedResponse) error {

			rows := pg_rsp.Rows()

			for rows.Next() {

				a, err := db.deriveAccountFromRows(rows)

				if err != nil {

					if !yield(nil, err) {
						return err
					}

					continue
				}

				if !yield(a, nil) {
					return nil
				}
			}

			err := rows.Close()

			if err != nil {
				return fmt.Errorf("Failed to iterate through database rows, %w", err)
			}

			return nil
		}

		pg_opts, err := countable.NewCountableOptions()

		if err != nil {
			yield(nil, fmt.Errorf("Failed to create pagination options, %w", err))
			return
		}

		q := fmt.Sprintf("SELECT id, account_type, name, display_name, blurb, url, public_key_uri, private_key_uri, created, lastmodified FROM %s", SQL_ACCOUNTS_TABLE_NAME)

		err = pg_sql.QueryPaginatedAll(db.database, pg_opts, pg_callback, q)

		if err != nil {
			yield(nil, fmt.Errorf("Failed to execute paginated query, %w", err))
			return
		}

	}
}

func (db *SQLAccountsDatabase) Close() error {
	return db.database.Close()
}

func (db *SQLAccountsDatabase) GetAccountIdsForDateRange(ctx context.Context, start int64, end int64) iter.Seq2[int64, error] {

	return func(yield func(int64, error) bool) {

		pg_callback := func(pg_rsp pg_sql.PaginatedResponse) error {

			rows := pg_rsp.Rows()

			for rows.Next() {

				var id int64

				err := rows.Scan(&id)

				if err != nil {

					if !yield(-1, err) {
						return err
					}

					continue
				}

				if !yield(id, nil) {
					return nil
				}
			}

			err := rows.Close()

			if err != nil {
				return fmt.Errorf("Failed to iterate through database rows, %w", err)
			}

			return nil
		}

		pg_opts, err := countable.NewCountableOptions()

		if err != nil {
			yield(-1, fmt.Errorf("Failed to create pagination options, %w", err))
		}

		q := fmt.Sprintf("SELECT id FROM %s WHERE created >= ? AND created <= ?", SQL_ACCOUNTS_TABLE_NAME)

		err = pg_sql.QueryPaginatedAll(db.database, pg_opts, pg_callback, q, start, end)

		if err != nil {
			yield(-1, fmt.Errorf("Failed to execute paginated query, %w", err))
		}
	}
}

func (db *SQLAccountsDatabase) GetAccountWithName(ctx context.Context, name string) (*activitypub.Account, error) {
	where := "name = ?"
	return db.getAccount(ctx, where, name)
}

func (db *SQLAccountsDatabase) getAccount(ctx context.Context, where string, args ...interface{}) (*activitypub.Account, error) {

	q := fmt.Sprintf("SELECT id, account_type, name, display_name, blurb, url, public_key_uri, private_key_uri, created, lastmodified FROM %s WHERE %s", SQL_ACCOUNTS_TABLE_NAME, where)

	row := db.database.QueryRowContext(ctx, q, args...)
	a, err := db.deriveAccountFromRows(row)

	if err != nil {
		return nil, err
	}

	return a, nil
}

func (db *SQLAccountsDatabase) deriveAccountFromRows(r any) (*activitypub.Account, error) {

	var id int64
	var account_type uint32
	var name string
	var display_name string
	var blurb string
	var url string
	var public_key_uri string
	var private_key_uri string
	var created int64
	var lastmod int64

	var err error

	switch r.(type) {
	case *sql.Row:
		err = r.(*sql.Row).Scan(&id, &account_type, &name, &display_name, &blurb, &url, &public_key_uri, &private_key_uri, &created, &lastmod)
	case *sql.Rows:
		err = r.(*sql.Rows).Scan(&id, &account_type, &name, &display_name, &blurb, &url, &public_key_uri, &private_key_uri, &created, &lastmod)
	default:
		return nil, fmt.Errorf("Invalid type %T", r)
	}

	switch {
	case err == sql.ErrNoRows:
		return nil, activitypub.ErrNotFound
	case err != nil:
		return nil, err
	default:
		//
	}

	a := &activitypub.Account{
		Id:            id,
		AccountType:   activitypub.AccountType(account_type),
		Name:          name,
		DisplayName:   display_name,
		Blurb:         blurb,
		URL:           url,
		PublicKeyURI:  public_key_uri,
		PrivateKeyURI: private_key_uri,
		Created:       created,
		LastModified:  lastmod,
	}

	return a, nil

}
