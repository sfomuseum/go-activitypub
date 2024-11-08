package database

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"

	pg_sql "github.com/aaronland/go-pagination-sql"
	"github.com/aaronland/go-pagination/countable"
	"github.com/sfomuseum/go-activitypub"
	"github.com/sfomuseum/go-activitypub/sqlite"
)

const SQL_ACCOUNTS_TABLE_NAME string = "accounts"

type SQLAccountsDatabase struct {
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

	db := &SQLAccountsDatabase{
		database: conn,
	}

	return db, nil
}

func (db *SQLAccountsDatabase) GetAccounts(ctx context.Context, acct GetAccountsCallbackFunc) error {
	return activitypub.ErrNotImplemented
}

func (db *SQLAccountsDatabase) GetAccountIdsForDateRange(ctx context.Context, start int64, end int64, cb GetAccountIdsCallbackFunc) error {

	pg_callback := func(pg_rsp pg_sql.PaginatedResponse) error {

		rows := pg_rsp.Rows()

		for rows.Next() {

			var id int64

			err := rows.Scan(&id)

			if err != nil {
				return fmt.Errorf("Failed to query database, %w", err)
			}

			err = cb(ctx, id)

			if err != nil {
				return fmt.Errorf("Failed to execute following callback for account %d, %w", id, err)
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

	q := fmt.Sprintf("SELECT id FROM %s WHERE created >= ? AND created <= ?", SQL_ACCOUNTS_TABLE_NAME)

	err = pg_sql.QueryPaginatedAll(db.database, pg_opts, pg_callback, q, start, end)

	if err != nil {
		return fmt.Errorf("Failed to execute paginated query, %w", err)
	}

	return nil
}

func (db *SQLAccountsDatabase) AddAccount(ctx context.Context, a *activitypub.Account) error {

	q := fmt.Sprintf("INSERT INTO %s (id, account_type, name, display_name, blurb, url, public_key_uri, private_key_uri, created, lastmodified) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", SQL_ACCOUNTS_TABLE_NAME)

	_, err := db.database.ExecContext(ctx, q, a.Id, a.AccountType, a.Name, a.DisplayName, a.Blurb, a.URL, a.PublicKeyURI, a.PrivateKeyURI, a.Created, a.LastModified)

	if err != nil {
		return fmt.Errorf("Failed to add account, %w", err)
	}

	return nil
}

func (db *SQLAccountsDatabase) GetAccountWithId(ctx context.Context, id int64) (*activitypub.Account, error) {
	where := "id = ?"
	return db.getAccount(ctx, where, id)
}

func (db *SQLAccountsDatabase) GetAccountWithName(ctx context.Context, name string) (*activitypub.Account, error) {
	where := "name = ?"
	return db.getAccount(ctx, where, name)
}

func (db *SQLAccountsDatabase) getAccount(ctx context.Context, where string, args ...interface{}) (*activitypub.Account, error) {

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

	q := fmt.Sprintf("SELECT id, account_type, name, display_name, blurb, url, public_key_uri, private_key_uri, created, lastmodified FROM %s WHERE %s", SQL_ACCOUNTS_TABLE_NAME, where)

	row := db.database.QueryRowContext(ctx, q, args...)

	err := row.Scan(&id, &account_type, &name, &display_name, &blurb, &url, &public_key_uri, &private_key_uri, &created, &lastmod)

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

func (db *SQLAccountsDatabase) UpdateAccount(ctx context.Context, acct *activitypub.Account) error {
	return activitypub.ErrNotImplemented
}

func (db *SQLAccountsDatabase) RemoveAccount(ctx context.Context, acct *activitypub.Account) error {
	return activitypub.ErrNotImplemented
}

func (db *SQLAccountsDatabase) Close(ctx context.Context) error {
	return db.database.Close()
}
