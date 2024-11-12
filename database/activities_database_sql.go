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

const SQL_ACTIVITIES_TABLE_NAME string = "activities"

type SQLActivitiesDatabase struct {
	ActivitiesDatabase
	database *sql.DB
}

func init() {
	ctx := context.Background()
	err := RegisterActivitiesDatabase(ctx, "sql", NewSQLActivitiesDatabase)

	if err != nil {
		panic(err)
	}
}

func NewSQLActivitiesDatabase(ctx context.Context, uri string) (ActivitiesDatabase, error) {

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

	db := &SQLActivitiesDatabase{
		database: conn,
	}

	return db, nil
}

func (db *SQLActivitiesDatabase) AddActivity(ctx context.Context, a *activitypub.Activity) error {

	q := fmt.Sprintf("INSERT INTO %s (id, activitypub_id, account_id, activity_type, activity_type_id, body, created) VALUES (?, ?, ?, ?, ?, ?, ?)", SQL_ACTIVITIES_TABLE_NAME)

	_, err := db.database.ExecContext(ctx, q, a.Id, a.ActivityPubId, a.AccountId, a.ActivityType, a.ActivityTypeId, a.Body, a.Created)

	if err != nil {
		return fmt.Errorf("Failed to add activity, %w", err)
	}

	return nil
}

func (db *SQLActivitiesDatabase) GetActivityWithId(ctx context.Context, id int64) (*activitypub.Activity, error) {

	where := "id = ?"
	return db.getActivity(ctx, where, id)
}

func (db *SQLActivitiesDatabase) GetActivityWithActivityPubId(ctx context.Context, id string) (*activitypub.Activity, error) {

	where := "activity_pub_id = ?"
	return db.getActivity(ctx, where, id)
}

func (db *SQLActivitiesDatabase) GetActivityWithActivityTypeAndId(ctx context.Context, activity_type activitypub.ActivityType, activity_type_id int64) (*activitypub.Activity, error) {

	where := "activity_type = ? AND activity_type_id = ?"
	return db.getActivity(ctx, where, activity_type, activity_type_id)
}

func (db *SQLActivitiesDatabase) GetActivities(ctx context.Context, cb GetActivitiesCallbackFunc) error {

	where := "1 = 1"
	args := make([]interface{}, 0)

	return db.getActivities(ctx, where, args, cb)
}

func (db *SQLActivitiesDatabase) GetActivitiesForAccount(ctx context.Context, id int64, cb GetActivitiesCallbackFunc) error {

	where := "account_id = ?"
	args := []interface{}{id}

	return db.getActivities(ctx, where, args, cb)
}

func (db *SQLActivitiesDatabase) Close(ctx context.Context) error {
	return db.database.Close()
}

func (db *SQLActivitiesDatabase) getActivity(ctx context.Context, where string, args ...interface{}) (*activitypub.Activity, error) {

	var id int64
	var activitypub_id string
	var account_id int64
	var activity_type int
	var activity_type_id int64
	var body string
	var created int64

	q := fmt.Sprintf("SELECT id, activitypub_id, account_id, activity_type, activity_type_id, body, created FROM %s WHERE %s", SQL_ACTIVITIES_TABLE_NAME, where)

	row := db.database.QueryRowContext(ctx, q, args...)

	err := row.Scan(&id, &activitypub_id, &account_id, &activity_type, &activity_type_id, &body, &created)

	switch {
	case err == sql.ErrNoRows:
		return nil, activitypub.ErrNotFound
	case err != nil:
		return nil, err
	default:
		//
	}

	a := &activitypub.Activity{
		Id:             id,
		ActivityPubId:  activitypub_id,
		AccountId:      account_id,
		ActivityType:   activitypub.ActivityType(activity_type),
		ActivityTypeId: activity_type_id,
		Body:           body,
		Created:        created,
	}

	return a, nil
}

func (db *SQLActivitiesDatabase) getActivities(ctx context.Context, where string, args []interface{}, cb GetActivitiesCallbackFunc) error {

	pg_callback := func(pg_rsp pg_sql.PaginatedResponse) error {

		rows := pg_rsp.Rows()

		for rows.Next() {

			var id int64
			var activitypub_id string
			var account_id int64
			var activity_type int
			var activity_type_id int64
			var body string
			var created int64

			err := rows.Scan(&id, &activitypub_id, &account_id, &activity_type, &activity_type_id, &body, &created)

			if err != nil {
				return fmt.Errorf("Failed to query database, %w", err)
			}

			a := &activitypub.Activity{
				Id:             id,
				ActivityPubId:  activitypub_id,
				AccountId:      account_id,
				ActivityType:   activitypub.ActivityType(activity_type),
				ActivityTypeId: activity_type_id,
				Body:           body,
				Created:        created,
			}

			err = cb(ctx, a)

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

	q := fmt.Sprintf("SELECT id, activitypub_id, account_id, activity_type, activity_type_id, body, created FROM %s WHERE %s", SQL_ACTIVITIES_TABLE_NAME, where)

	err = pg_sql.QueryPaginatedAll(db.database, pg_opts, pg_callback, q, args...)

	if err != nil {
		return fmt.Errorf("Failed to execute paginated query, %w", err)
	}

	return nil

}
