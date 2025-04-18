package followers

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/sfomuseum/go-activitypub"
	"github.com/sfomuseum/go-activitypub/ap"
	"github.com/sfomuseum/go-activitypub/database"
	"github.com/sfomuseum/go-activitypub/uris"
)

func CountFollowers(ctx context.Context, db database.FollowersDatabase, account_id int64) (uint32, error) {

	count := uint32(0)

	followers_cb := func(ctx context.Context, follower string) error {
		atomic.AddUint32(&count, 1)
		return nil
	}

	err := db.GetFollowersForAccount(ctx, account_id, followers_cb)

	if err != nil {
		return 0, fmt.Errorf("Failed to count followers, %w", err)
	}

	return atomic.LoadUint32(&count), nil
}

func GetFollower(ctx context.Context, db database.FollowersDatabase, account_id int64, follower_address string) (*activitypub.Follower, error) {
	return db.GetFollower(ctx, account_id, follower_address)
}

func AddFollower(ctx context.Context, db database.FollowersDatabase, account_id int64, follower_address string) (int64, error) {

	f, err := activitypub.NewFollower(ctx, account_id, follower_address)

	if err != nil {
		return -1, fmt.Errorf("Failed to create new follower, %w", err)
	}

	err = db.AddFollower(ctx, f)

	if err != nil {
		return -1, fmt.Errorf("Failed to add follower, %w", err)
	}

	return f.Id, nil
}

// Is follower_address following account_id?
func IsFollower(ctx context.Context, db database.FollowersDatabase, account_id int64, follower_address string) (bool, *activitypub.Follower, error) {

	f, err := GetFollower(ctx, db, account_id, follower_address)

	if err == nil {
		return true, f, nil
	}

	if err == activitypub.ErrNotFound {
		return false, nil, nil
	}

	return false, nil, fmt.Errorf("Failed to follower record, %w", err)
}

func FollowersResource(ctx context.Context, uris_table *uris.URIs, followers_database database.FollowersDatabase, a *activitypub.Account) (*ap.Followers, error) {

	followers_path := uris.AssignResource(uris_table.Followers, a.Name)
	followers_url := uris.NewURL(uris_table, followers_path)

	count, err := CountFollowers(ctx, followers_database, a.Id)

	if err != nil {
		return nil, fmt.Errorf("Failed to count followers, %w", err)
	}

	f := &ap.Followers{
		Context:    ap.ACTIVITYSTREAMS_CONTEXT,
		Id:         followers_url.String(),
		Type:       "OrderedCollection",
		TotalItems: count,
		First:      followers_url.String(),
	}

	return f, nil
}
