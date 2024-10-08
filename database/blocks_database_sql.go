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

const SQL_BLOCKS_TABLE_NAME string = "blocks"

type SQLBlocksDatabase struct {
	BlocksDatabase
	database *sql.DB
}

func init() {
	ctx := context.Background()
	RegisterBlocksDatabase(ctx, "sql", NewSQLBlocksDatabase)
}

func NewSQLBlocksDatabase(ctx context.Context, uri string) (BlocksDatabase, error) {

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
			return nil, fmt.Errorf("Failed to configure SQLite connection, %w", err)
		}
	}

	db := &SQLBlocksDatabase{
		database: conn,
	}

	return db, nil
}

func (db *SQLBlocksDatabase) GetBlockIdsForDateRange(ctx context.Context, start int64, end int64, cb GetBlockIdsCallbackFunc) error {

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
				return fmt.Errorf("Failed to execute following callback for block %d, %w", id, err)
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

	q := fmt.Sprintf("SELECT id FROM %s WHERE created >= ? AND created <= ?", SQL_BLOCKS_TABLE_NAME)

	err = pg_sql.QueryPaginatedAll(db.database, pg_opts, pg_callback, q, start, end)

	if err != nil {
		return fmt.Errorf("Failed to execute paginated query, %w", err)
	}

	return nil
}

func (db *SQLBlocksDatabase) GetBlockWithId(ctx context.Context, block_id int64) (*Block, error) {
	where := "id = ?"
	return db.getBlock(ctx, where, block_id)
}

func (db *SQLBlocksDatabase) GetBlockWithAccountIdAndAddress(ctx context.Context, account_id int64, host string, name string) (*Block, error) {

	where := "account_id = ? AND host = ? AND name = ?"
	return db.getBlock(ctx, where, account_id, host, name)

}

func (db *SQLBlocksDatabase) getBlock(ctx context.Context, where string, args ...interface{}) (*Block, error) {

	q := fmt.Sprintf("SELECT id, account_id, name, host, created, lastmodified FROM %s WHERE %s", SQL_BLOCKS_TABLE_NAME, where)
	row := db.database.QueryRowContext(ctx, q, args...)

	var id int64
	var account_id int64
	var name string
	var host string
	var created int64
	var lastmod int64

	err := row.Scan(&id, &account_id, &name, &host, &created, &lastmod)

	switch {
	case err == sql.ErrNoRows:
		return nil, ErrNotFound
	case err != nil:
		return nil, fmt.Errorf("Failed to query database, %w", err)
	default:

		n := &Block{
			Id:           id,
			AccountId:    account_id,
			Name:         name,
			Host:         host,
			Created:      created,
			LastModified: lastmod,
		}

		return n, nil
	}

}

func (db *SQLBlocksDatabase) AddBlock(ctx context.Context, block *Block) error {

	q := fmt.Sprintf("INSERT INTO %s (id, account_id, name, host, created, lastmodified) VALUES (?, ?, ?, ?, ?, ?)", SQL_BLOCKS_TABLE_NAME)

	_, err := db.database.ExecContext(ctx, q, block.Id, block.AccountId, block.Name, block.Host, block.Created, block.LastModified)

	if err != nil {
		return fmt.Errorf("Failed to add block, %w", err)
	}

	return nil
}

func (db *SQLBlocksDatabase) UpdateBlock(ctx context.Context, block *Block) error {

	q := fmt.Sprintf("UPDATE %s SET account_id=?, name=?, host=?, created=?, lastmodified=? WHERE id = ?", SQL_BLOCKS_TABLE_NAME)

	_, err := db.database.ExecContext(ctx, q, block.AccountId, block.Name, block.Host, block.Created, block.LastModified, block.Id)

	if err != nil {
		return fmt.Errorf("Failed to add block, %w", err)
	}

	return nil
}

func (db *SQLBlocksDatabase) RemoveBlock(ctx context.Context, block *Block) error {

	q := fmt.Sprintf("DELETE FROM %s WHERE id = ?", SQL_BLOCKS_TABLE_NAME)

	_, err := db.database.ExecContext(ctx, q, block.Id)

	if err != nil {
		return fmt.Errorf("Failed to remove block, %w", err)
	}

	return nil
}

func (db *SQLBlocksDatabase) Close(ctx context.Context) error {
	return db.database.Close()
}
