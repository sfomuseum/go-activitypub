package activitypub

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"

	"github.com/sfomuseum/go-activitypub/sqlite"
)

const SQL_ACCOUNTS_TABLE_NAME string = "accounts"

type SQLAccountsDatabase struct {
	AccountsDatabase
	database *sql.DB
}

func init() {
	ctx := context.Background()
	RegisterAccountsDatabase(ctx, "sql", NewSQLAccountsDatabase)
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

func (db *SQLAccountsDatabase) AddAccount(ctx context.Context, a *Account) error {

	q := fmt.Sprintf("INSERT INTO %s (id, account_type, name, display_name, blurb, url, public_key_uri, private_key_uri, created, lastmodified) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", SQL_ACCOUNTS_TABLE_NAME)

	_, err := db.database.ExecContext(ctx, q, a.Id, a.AccountType, a.Name, a.DisplayName, a.Blurb, a.URL, a.PublicKeyURI, a.PrivateKeyURI, a.Created, a.LastModified)

	if err != nil {
		return fmt.Errorf("Failed to add account, %w", err)
	}

	return nil
}

func (db *SQLAccountsDatabase) GetAccountWithId(ctx context.Context, id int64) (*Account, error) {
	where := "id = ?"
	return db.getAccount(ctx, where, id)
}

func (db *SQLAccountsDatabase) GetAccountWithName(ctx context.Context, name string) (*Account, error) {
	where := "name = ?"
	return db.getAccount(ctx, where, name)
}

func (db *SQLAccountsDatabase) getAccount(ctx context.Context, where string, args ...interface{}) (*Account, error) {

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
		return nil, ErrNotFound
	case err != nil:
		return nil, err
	default:
		//
	}

	a := &Account{
		Id:            id,
		AccountType:   AccountType(account_type),
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

func (db *SQLAccountsDatabase) Close(ctx context.Context) error {
	return db.database.Close()
}
