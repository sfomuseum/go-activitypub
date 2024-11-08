package database

import (
	"context"
	"fmt"
	"io"

	_ "log/slog"

	aa_docstore "github.com/aaronland/gocloud-docstore"
	"github.com/sfomuseum/go-activitypub"
	gc_docstore "gocloud.dev/docstore"
)

type DocstoreFollowingDatabase struct {
	FollowingDatabase
	collection *gc_docstore.Collection
}

func init() {

	ctx := context.Background()

	err := RegisterFollowingDatabase(ctx, "awsdynamodb", NewDocstoreFollowingDatabase)

	if err != nil {
		panic(err)
	}

	for _, scheme := range gc_docstore.DefaultURLMux().CollectionSchemes() {
		err := RegisterFollowingDatabase(ctx, scheme, NewDocstoreFollowingDatabase)

		if err != nil {
			panic(err)
		}
	}

}

func NewDocstoreFollowingDatabase(ctx context.Context, uri string) (FollowingDatabase, error) {

	col, err := aa_docstore.OpenCollection(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to open collection, %w", err)
	}

	db := &DocstoreFollowingDatabase{
		collection: col,
	}

	return db, nil
}

func (db *DocstoreFollowingDatabase) GetFollowingIdsForDateRange(ctx context.Context, start int64, end int64, cb GetFollowingIdsCallbackFunc) error {

	q := db.collection.Query()
	q = q.Where("Created", ">=", start)
	q = q.Where("Created", "<=", end)

	iter := q.Get(ctx, "Id")
	defer iter.Stop()

	for {

		var f activitypub.Following
		err := iter.Next(ctx, &f)

		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("Failed to interate, %w", err)
		} else {
			err := cb(ctx, f.Id)

			if err != nil {
				return fmt.Errorf("Failed to invoke callback for following %d, %w", f.Id, err)
			}
		}
	}

	return nil
}

func (db *DocstoreFollowingDatabase) GetFollowing(ctx context.Context, account_id int64, following_address string) (*activitypub.Following, error) {

	q := db.collection.Query()
	q = q.Where("AccountId", "=", account_id)
	q = q.Where("FollowingAddress", "=", following_address)

	iter := q.Get(ctx)
	defer iter.Stop()

	for {

		var f activitypub.Following
		err := iter.Next(ctx, &f)

		if err == io.EOF {
			break
		} else if err != nil {
			return nil, fmt.Errorf("Failed to interate, %w", err)
		} else {
			return &f, nil
		}
	}

	return nil, activitypub.ErrNotFound
}

func (db *DocstoreFollowingDatabase) AddFollowing(ctx context.Context, f *activitypub.Following) error {

	return db.collection.Put(ctx, f)
}

func (db *DocstoreFollowingDatabase) RemoveFollowing(ctx context.Context, f *activitypub.Following) error {

	return db.collection.Delete(ctx, f)
}

func (db *DocstoreFollowingDatabase) GetFollowingForAccount(ctx context.Context, account_id int64, following_callback GetFollowingCallbackFunc) error {

	q := db.collection.Query()
	q = q.Where("AccountId", "=", account_id)

	iter := q.Get(ctx)
	defer iter.Stop()

	for {

		var f activitypub.Following
		err := iter.Next(ctx, &f)

		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("Failed to interate, %w", err)
		} else {

			err := following_callback(ctx, f.FollowingAddress)

			if err != nil {
				return fmt.Errorf("Failed to execute following callback for '%s', %w", f.FollowingAddress, err)
			}
		}
	}

	return nil
}

func (db *DocstoreFollowingDatabase) Close(ctx context.Context) error {
	return db.collection.Close()
}
