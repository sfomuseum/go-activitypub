package activitypub

import (
	"context"
	"fmt"
	"io"

	_ "log/slog"

	aa_docstore "github.com/aaronland/gocloud-docstore"
	gc_docstore "gocloud.dev/docstore"
)

type DocstoreFollowingDatabase struct {
	FollowingDatabase
	collection *gc_docstore.Collection
}

func init() {

	ctx := context.Background()

	RegisterFollowingDatabase(ctx, "awsdynamodb", NewDocstoreFollowingDatabase)

	for _, scheme := range gc_docstore.DefaultURLMux().CollectionSchemes() {
		RegisterFollowingDatabase(ctx, scheme, NewDocstoreFollowingDatabase)
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

func (db *DocstoreFollowingDatabase) GetFollowing(ctx context.Context, account_id int64, following_address string) (*Following, error) {

	q := db.collection.Query()
	q = q.Where("AccountId", "=", account_id)
	q = q.Where("FollowingAddress", "=", following_address)

	iter := q.Get(ctx)
	defer iter.Stop()

	for {

		var f Following
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

func (db *DocstoreFollowingDatabase) AddFollowing(ctx context.Context, f *Following) error {

	return db.collection.Put(ctx, f)
}

func (db *DocstoreFollowingDatabase) RemoveFollowing(ctx context.Context, f *Following) error {

	return db.collection.Delete(ctx, f)
}

func (db *DocstoreFollowingDatabase) GetFollowingForAccount(ctx context.Context, account_id int64, following_callback GetFollowingCallbackFunc) error {

	q := db.collection.Query()
	q = q.Where("AccountId", "=", account_id)

	iter := q.Get(ctx)
	defer iter.Stop()

	for {

		var f Following
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
