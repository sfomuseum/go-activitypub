package activitypub

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
)

type SQLActorDatabase struct {
	ActorDatabase
	database *sql.DB
}

func init() {
	ctx := context.Background()
	RegisterActorDatabase(ctx, "sql", NewSQLActorDatabase)
}

func NewSQLActorDatabase(ctx context.Context, uri string) (ActorDatabase, error) {

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

	db := &SQLActorDatabase{
		database: conn,
	}

	return db, nil
}

func (db *SQLActorDatabase) AddActor(ctx context.Context, a *Actor) error {

	return nil
}
