package activitypub

import (
	"context"
	"fmt"
	"io"

	aa_docstore "github.com/aaronland/gocloud-docstore"
	gc_docstore "gocloud.dev/docstore"
)

type DocstoreFollowersDatabase struct {
	FollowersDatabase
	collection *gc_docstore.Collection
}

func init() {

	ctx := context.Background()

	RegisterFollowersDatabase(ctx, "awsdynamodb", NewDocstoreFollowersDatabase)

	for _, scheme := range gc_docstore.DefaultURLMux().CollectionSchemes() {
		RegisterFollowersDatabase(ctx, scheme, NewDocstoreFollowersDatabase)
	}

}

func NewDocstoreFollowersDatabase(ctx context.Context, uri string) (FollowersDatabase, error) {

	col, err := aa_docstore.OpenCollection(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to open collection, %w", err)
	}

	db := &DocstoreFollowersDatabase{
		collection: col,
	}

	return db, nil
}

func (db *DocstoreFollowersDatabase) GetFollowerIdsForDateRange(ctx context.Context, start int64, end int64, cb GetFollowerIdsCallbackFunc) error {

	q := db.collection.Query()
	q = q.Where("Created", ">=", start)
	q = q.Where("Created", "<=", end)

	return db.getFollowerIdsWithCallback(ctx, q, cb)
}

func (db *DocstoreFollowersDatabase) GetAllFollowers(ctx context.Context, cb GetFollowersCallbackFunc) error {

	q := db.collection.Query()
	return db.getFollowerAddressesWithCallback(ctx, q, cb)
}

func (db *DocstoreFollowersDatabase) HasFollowers(ctx context.Context, account_id int64) (bool, error) {

	q := db.collection.Query()
	q = q.Where("AccountId", "=", account_id)
	q = q.Limit(1)

	iter := q.Get(ctx, "Id")
	defer iter.Stop()

	var f Follower
	err := iter.Next(ctx, &f)

	if err == io.EOF {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("Failed to interate, %w", err)
	} else {
		return true, nil
	}

}

func (db *DocstoreFollowersDatabase) GetFollower(ctx context.Context, account_id int64, follower_address string) (*Follower, error) {

	q := db.collection.Query()
	q = q.Where("AccountId", "=", account_id)
	q = q.Where("FollowerAddress", "=", follower_address)

	iter := q.Get(ctx)
	defer iter.Stop()

	for {

		var f Follower
		err := iter.Next(ctx, &f)

		if err == io.EOF {
			break
		} else if err != nil {
			return nil, fmt.Errorf("Failed to interate, %w", err)
		} else {
			return &f, nil
		}
	}

	return nil, ErrNotFound
}

func (db *DocstoreFollowersDatabase) AddFollower(ctx context.Context, f *Follower) error {

	return db.collection.Put(ctx, f)
}

func (db *DocstoreFollowersDatabase) RemoveFollower(ctx context.Context, f *Follower) error {

	return db.collection.Delete(ctx, f)
}

func (db *DocstoreFollowersDatabase) GetFollowersForAccount(ctx context.Context, account_id int64, cb GetFollowersCallbackFunc) error {

	q := db.collection.Query()
	q = q.Where("AccountId", "=", account_id)

	return db.getFollowerAddressesWithCallback(ctx, q, cb)
}

func (db *DocstoreFollowersDatabase) Close(ctx context.Context) error {
	return db.collection.Close()
}

func (db *DocstoreFollowersDatabase) getFollowerIdsWithCallback(ctx context.Context, q *gc_docstore.Query, cb GetFollowerIdsCallbackFunc) error {

	iter := q.Get(ctx, "Id")
	defer iter.Stop()

	for {

		var f Follower
		err := iter.Next(ctx, &f)

		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("Failed to interate, %w", err)
		} else {
			err := cb(ctx, f.Id)

			if err != nil {
				return fmt.Errorf("Failed to invoke callback for follower %d, %w", f.Id, err)
			}
		}
	}

	return nil
}

func (db *DocstoreFollowersDatabase) getFollowerAddressesWithCallback(ctx context.Context, q *gc_docstore.Query, cb GetFollowersCallbackFunc) error {

	iter := q.Get(ctx, "FollowerAddress")
	defer iter.Stop()

	for {

		var f Follower
		err := iter.Next(ctx, &f)

		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("Failed to interate, %w", err)
		} else {

			err := cb(ctx, f.FollowerAddress)

			if err != nil {
				return fmt.Errorf("Failed to execute followers callback for '%s', %w", f.FollowerAddress, err)
			}
		}
	}

	return nil
}
