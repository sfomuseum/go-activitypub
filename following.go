package activitypub

import (
	"context"
	"fmt"
	"log/slog"
	"sync/atomic"
	"time"

	"github.com/sfomuseum/go-activitypub/id"
)

type Following struct {
	Id               int64  `json:"id"`
	AccountId        int64  `json:"account_id"`
	FollowingAddress string `json:"following_address"`
	Created          int64  `json:"created"`
}

func CountFollowing(ctx context.Context, db FollowingDatabase, account_id int64) (uint32, error) {

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

func GetFollowing(ctx context.Context, db FollowingDatabase, account_id int64, following_address string) (*Following, error) {

	slog.Debug("Get following", "account", account_id, "following", following_address)

	return db.GetFollowing(ctx, account_id, following_address)
}

func AddFollowing(ctx context.Context, db FollowingDatabase, account_id int64, following_address string) error {

	slog.Debug("Add following", "account", account_id, "following", following_address)

	f, err := NewFollowing(ctx, account_id, following_address)

	if err != nil {
		return fmt.Errorf("Failed to create new following, %w", err)
	}

	return db.AddFollowing(ctx, f)
}

func NewFollowing(ctx context.Context, account_id int64, following_address string) (*Following, error) {

	db_id, err := id.NewId()

	if err != nil {
		return nil, fmt.Errorf("Failed to create new following ID, %w", err)
	}

	now := time.Now()
	ts := now.Unix()

	b := &Following{
		Id:               db_id,
		AccountId:        account_id,
		FollowingAddress: following_address,
		Created:          ts,
	}

	return b, nil
}

// Is account_id following following_address?
func IsFollowing(ctx context.Context, db FollowingDatabase, account_id int64, following_address string) (bool, *Following, error) {

	f, err := GetFollowing(ctx, db, account_id, following_address)

	if err == nil {
		return true, f, nil
	}

	if err == ErrNotFound {
		return false, nil, nil
	}

	return false, nil, fmt.Errorf("Failed to following record, %w", err)
}
