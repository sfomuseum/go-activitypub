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

const SQL_BOOSTS_TABLE_NAME string = "boosts"

type SQLBoostsDatabase struct {
	BoostsDatabase
	database *sql.DB
}

func init() {

	ctx := context.Background()
	err := RegisterBoostsDatabase(ctx, "sql", NewSQLBoostsDatabase)

	if err != nil {
		panic(err)
	}
}

func NewSQLBoostsDatabase(ctx context.Context, uri string) (BoostsDatabase, error) {

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
			return nil, fmt.Errorf("Failed to configure SQLite, %w", err)
		}
	}

	db := &SQLBoostsDatabase{
		database: conn,
	}

	return db, nil
}

func (db *SQLBoostsDatabase) GetBoostWithId(ctx context.Context, id int64) (*activitypub.Boost, error) {

	where := "id = ?"
	return db.getBoost(ctx, where, id)
}

func (db *SQLBoostsDatabase) GetBoostWithPostIdAndActor(ctx context.Context, post_id int64, actor string) (*activitypub.Boost, error) {

	where := "post_id = ? AND actor = ?"
	return db.getBoost(ctx, where, post_id, actor)
}

func (db *SQLBoostsDatabase) GetBoostsForAccount(ctx context.Context, account_id int64, cb GetBoostsCallbackFunc) error {

	where := "account_id = ?"

	args := []interface{}{
		account_id,
	}

	return db.getBoostsForQuery(ctx, where, args, cb)
}

func (db *SQLBoostsDatabase) GetBoostsForPosts(ctx context.Context, post_id int64, cb GetBoostsCallbackFunc) error {

	where := "post_id = ?"

	args := []interface{}{
		post_id,
	}

	return db.getBoostsForQuery(ctx, where, args, cb)
}

func (db *SQLBoostsDatabase) GetBoostsForPostIdAndActor(ctx context.Context, post_id int64, actor string, cb GetBoostsCallbackFunc) error {

	where := "post_id = ? AND actor = ?"
	args := []interface{}{
		post_id,
		actor,
	}

	return db.getBoostsForQuery(ctx, where, args, cb)
}

func (db *SQLBoostsDatabase) AddBoost(ctx context.Context, b *activitypub.Boost) error {

	q := fmt.Sprintf("INSERT INTO %s (id, account_id, psot_id, actor, created) VALUES (?, ?, ?, ?, ?)", SQL_BOOSTS_TABLE_NAME)

	_, err := db.database.ExecContext(ctx, q, b.Id, b.AccountId, b.PostId, b.Actor, b.Created)

	if err != nil {
		return fmt.Errorf("Failed to add boost, %w", err)
	}

	return nil
}

func (db *SQLBoostsDatabase) RemoveBoost(ctx context.Context, b *activitypub.Boost) error {

	q := fmt.Sprintf("DELETE FROM %s WHERE id= ?", SQL_BOOSTS_TABLE_NAME)

	_, err := db.database.ExecContext(ctx, q, b.Id)

	if err != nil {
		return fmt.Errorf("Failed to remove boost, %w", err)
	}

	return nil
}

func (db *SQLBoostsDatabase) Close(ctx context.Context) error {
	return db.database.Close()
}

func (db *SQLBoostsDatabase) getBoost(ctx context.Context, where string, args ...interface{}) (*activitypub.Boost, error) {

	var id int64
	var account_id int64
	var post_id int64
	var actor string
	var created int64

	q := fmt.Sprintf("SELECT id, account_id, post_id, actor, created FROM %s WHERE %s", SQL_BOOSTS_TABLE_NAME, where)

	row := db.database.QueryRowContext(ctx, q, args...)

	err := row.Scan(&id, &account_id, &post_id, &actor, &created)

	switch {
	case err == sql.ErrNoRows:
		return nil, activitypub.ErrNotFound
	case err != nil:
		return nil, err
	default:
		//
	}

	b := &activitypub.Boost{
		Id:        id,
		AccountId: account_id,
		PostId:    post_id,
		Actor:     actor,
		Created:   created,
	}

	return b, nil
}

func (db *SQLBoostsDatabase) getBoostsForQuery(ctx context.Context, where string, args []interface{}, cb GetBoostsCallbackFunc) error {

	pg_callback := func(pg_rsp pg_sql.PaginatedResponse) error {

		rows := pg_rsp.Rows()

		for rows.Next() {

			var id int64
			var account_id int64
			var post_id int64
			var actor string
			var created int64

			err := rows.Scan(&id, &account_id, &post_id, &actor, &created)

			switch {
			case err == sql.ErrNoRows:
				return nil
			case err != nil:
				return err
			default:

				b := &activitypub.Boost{
					Id:        id,
					AccountId: account_id,
					PostId:    post_id,
					Actor:     actor,
					Created:   created,
				}

				err = cb(ctx, b)

				if err != nil {
					return fmt.Errorf("Failed to execute callback for '%b', %w", b.Id, err)
				}

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

	q := fmt.Sprintf("SELECT id, account_id, post_id, actor, created FROM %s WHERE %s", SQL_BOOSTS_TABLE_NAME, where)

	err = pg_sql.QueryPaginatedAll(db.database, pg_opts, pg_callback, q, args...)

	if err != nil {
		return fmt.Errorf("Failed to execute paginated query, %w", err)
	}

	return nil
}
