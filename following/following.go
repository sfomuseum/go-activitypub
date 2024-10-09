package following

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/sfomuseum/go-activitypub"
	"github.com/sfomuseum/go-activitypub/ap"
	"github.com/sfomuseum/go-activitypub/database"
	"github.com/sfomuseum/go-activitypub/uris"
)

func CountFollowing(ctx context.Context, db database.FollowingDatabase, account_id int64) (uint32, error) {

	count := uint32(0)

	following_cb := func(ctx context.Context, following string) error {
		atomic.AddUint32(&count, 1)
		return nil
	}

	err := db.GetFollowingForAccount(ctx, account_id, following_cb)

	if err != nil {
		return 0, fmt.Errorf("Failed to count following, %w", err)
	}

	return atomic.LoadUint32(&count), nil
}

func GetFollowing(ctx context.Context, db database.FollowingDatabase, account_id int64, following_address string) (*activitypub.Following, error) {

	return db.GetFollowing(ctx, account_id, following_address)
}

func AddFollowing(ctx context.Context, db database.FollowingDatabase, account_id int64, following_address string) error {

	f, err := activitypub.NewFollowing(ctx, account_id, following_address)

	if err != nil {
		return fmt.Errorf("Failed to create new following, %w", err)
	}

	return db.AddFollowing(ctx, f)
}

// Is account_id following following_address?
func IsFollowing(ctx context.Context, db database.FollowingDatabase, account_id int64, following_address string) (bool, *activitypub.Following, error) {

	f, err := GetFollowing(ctx, db, account_id, following_address)

	if err == nil {
		return true, f, nil
	}

	if err == activitypub.ErrNotFound {
		return false, nil, nil
	}

	return false, nil, fmt.Errorf("Failed to following record, %w", err)
}

func FollowingResource(ctx context.Context, uris_table *uris.URIs, following_database database.FollowingDatabase, a *activitypub.Account) (*ap.Following, error) {

	following_path := uris.AssignResource(uris_table.Following, a.Name)
	following_url := uris.NewURL(uris_table, following_path)

	count, err := CountFollowing(ctx, following_database, a.Id)

	if err != nil {
		return nil, fmt.Errorf("Failed to count following, %w", err)
	}

	f := &ap.Following{
		Context:    "https://www.w3.org/ns/activitystreams",
		Id:         following_url.String(),
		Type:       "OrderedCollection",
		TotalItems: count,
		First:      following_url.String(),
	}

	return f, nil
}
