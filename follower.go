package activitypub

import (
	"context"
	"fmt"
	"time"
)

type Follower struct {
	Id              int64  `json:"id"`
	AccountId       int64  `json:"account_id"`
	FollowerAddress string `json:"follower_address"`
	Created         int64  `json:"created"`
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
		FollowerAddress: follower_address,
		Created:         ts,
	}

	return b, nil
}

func IsFollowingAccount(ctx context.Context, db FollowersDatabase, follower_address string, account_id int64) (bool, *Follower, error) {

	f, err := db.GetFollower(ctx, account_id, follower_address)

	if err == nil {
		return true, f, nil
	}

	if err == ErrNotFound {
		return false, nil, nil
	}

	return false, nil, fmt.Errorf("Failed to follower record, %w", err)
}
