package activitypub

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

type Follower struct {
	Id              int64  `json:"id"`
	AccountId       int64  `json:"account_id"`
	FollowerAddress string `json:"follower_address"`
	Created         int64  `json:"created"`
}

func GetFollower(ctx context.Context, db FollowersDatabase, account_id int64, follower_address string) (*Follower, error) {

	slog.Info("Get follower", "account", account_id, "follower", follower_address)

	return db.GetFollower(ctx, account_id, follower_address)
}

func AddFollower(ctx context.Context, db FollowersDatabase, account_id int64, follower_address string) error {

	slog.Info("Add follower", "account", account_id, "follower", follower_address)

	f, err := NewFollower(ctx, account_id, follower_address)

	if err != nil {
		return fmt.Errorf("Failed to create new follower, %w", err)
	}

	return db.AddFollower(ctx, f)
}

func NewFollower(ctx context.Context, account_id int64, follower_address string) (*Follower, error) {

	db_id, err := NewId()

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
