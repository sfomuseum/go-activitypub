package activitypub

import (
	"context"
	"fmt"
	"io"
	"time"

	aa_docstore "github.com/aaronland/gocloud-docstore"
	gc_docstore "gocloud.dev/docstore"
)

type DocstoreFollowersDatabase struct {
	FollowersDatabase
	collection *gc_docstore.Collection
}

func init() {
	// ctx := context.Background()
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

func (db *DocstoreFollowersDatabase) IsFollowing(ctx context.Context, follower_address string, account_id int64) (bool, error) {

	_, err := db.getFollower(ctx, follower_address, account_id)

	switch {
	case err == ErrNotFound:
		return false, nil
	case err != nil:
		return false, err
	default:
		return true, nil
	}
}

func (db *DocstoreFollowersDatabase) getFollower(ctx context.Context, follower_address string, account_id int64) (*Follower, error) {

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

func (db *DocstoreFollowersDatabase) AddFollower(ctx context.Context, account_id int64, follower_address string) error {

	now := time.Now()
	ts := now.Unix()

	f := &Follower{
		AccountId:       account_id,
		FollowerAddress: follower_address,
		Created:         ts,
	}

	return db.collection.Put(ctx, f)
}

func (db *DocstoreFollowersDatabase) RemoveFollower(ctx context.Context, account_id int64, follower_address string) error {

	f, err := db.getFollower(ctx, follower_address, account_id)

	switch {
	case err == ErrNotFound:
		return nil
	case err != nil:
		return err
	default:
		return db.collection.Delete(ctx, f)
	}
}

func (db *DocstoreFollowersDatabase) GetFollowers(ctx context.Context, account_id int64, followers_callback GetFollowersCallbackFunc) error {

	q := db.collection.Query()
	q = q.Where("AccountId", "=", account_id)

	iter := q.Get(ctx)
	defer iter.Stop()

	for {

		var f Follower
		err := iter.Next(ctx, &f)

		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("Failed to interate, %w", err)
		} else {

			err := followers_callback(ctx, f.FollowerAddress)

			if err != nil {
				return fmt.Errorf("Failed to execute followers callback for '%s', %w", f.FollowerAddress, err)
			}
		}
	}

	return nil
}