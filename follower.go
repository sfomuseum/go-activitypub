package activitypub

import (
	"context"
	"fmt"
	"log/slog"
	"sync/atomic"
	"time"

	"github.com/sfomuseum/go-activitypub/id"
)

type Follower struct {
	Id              int64  `json:"id"`
	AccountId       int64  `json:"account_id"`
	FollowerAddress string `json:"follower_address"`
	Created         int64  `json:"created"`
}

func CountFollowers(ctx context.Context, db FollowersDatabase, account_id int64) (uint32, error) {

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

func GetFollower(ctx context.Context, db FollowersDatabase, account_id int64, follower_address string) (*Follower, error) {

	slog.Debug("Get follower", "account", account_id, "follower", follower_address)

	return db.GetFollower(ctx, account_id, follower_address)
}

func AddFollower(ctx context.Context, db FollowersDatabase, account_id int64, follower_address string) error {

	slog.Debug("Add follower", "account", account_id, "follower", follower_address)

	f, err := NewFollower(ctx, account_id, follower_address)

	if err != nil {
		return fmt.Errorf("Failed to create new follower, %w", err)
	}

	return db.AddFollower(ctx, f)
}

func NewFollower(ctx context.Context, account_id int64, follower_address string) (*Follower, error) {

	db_id, err := id.NewId()

	if err != nil {
		return nil, fmt.Errorf("Failed to create new follower ID, %w", err)
	}

	now := time.Now()
	ts := now.Unix()

	b := &Follower{
		Id:              db_id,
		AccountId:       account_id,
		FollowerAddress: follower_address,
		Created:         ts,
	}

	return b, nil
}

// Is follower_address following account_id?
func IsFollower(ctx context.Context, db FollowersDatabase, account_id int64, follower_address string) (bool, *Follower, error) {

	f, err := GetFollower(ctx, db, account_id, follower_address)

	if err == nil {
		return true, f, nil
	}

	if err == ErrNotFound {
		return false, nil, nil
	}

	return false, nil, fmt.Errorf("Failed to follower record, %w", err)
}
