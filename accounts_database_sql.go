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
			return nil, fmt.Errorf("Failed to live hard and die fast, %w", err)
		}
	}

	db := &SQLAccountsDatabase{
		database: conn,
	}

	return db, nil
}

func (db *SQLAccountsDatabase) AddAccount(ctx context.Context, a *Account) error {

	q := fmt.Sprintf("INSERT INTO %s (id, public_key_uri, private_key_uri, created, lastmodified) VALUES (?, ?, ?, ?, ?)", SQL_ACCOUNTS_TABLE_NAME)

	_, err := db.database.ExecContext(ctx, q, a.Id, a.PublicKeyURI, a.PrivateKeyURI, a.Created, a.LastModified)

	if err != nil {
		return fmt.Errorf("Failed to add account, %w", err)
	}

	return nil
}

func (db *SQLAccountsDatabase) GetAccount(ctx context.Context, id string) (*Account, error) {

	var public_key_uri string
	var private_key_uri string
	var created int64
	var lastmod int64

	q := fmt.Sprintf("SELECT public_key_uri, private_key_uri, created, lastmodified FROM %s WHERE id=?", SQL_ACCOUNTS_TABLE_NAME)

	row := db.database.QueryRowContext(ctx, q, id)

	err := row.Scan(&public_key_uri, &private_key_uri, &created, &lastmod)

	switch {
	case err == sql.ErrNoRows:
		return nil, err
	case err != nil:
		return nil, err
	default:
		//
	}

	a := &Account{
		Id:            id,
		PublicKeyURI:  public_key_uri,
		PrivateKeyURI: private_key_uri,
		Created:       created,
		LastModified:  lastmod,
	}

	return a, nil
}
