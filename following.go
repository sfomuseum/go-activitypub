package activitypub

import (
	"context"
	"fmt"
	"time"
)

type Following struct {
	Id               int64  `json:"id"`
	AccountId        int64  `json:"account_id"`
	FollowingAddress string `json:"following_address"`
	Created          int64  `json:"created"`
}

func NewFollowing(ctx context.Context, account_id int64, following_address string) (*Following, error) {

	db_id, err := NewId()

	if err != nil {
		return nil, fmt.Errorf("Failed to create new following ID, %w", err)
	}

	now := time.Now()
	ts := now.Unix()

	b := &Following{
		Id:               db_id,
		FollowingAddress: following_address,
		Created:          ts,
	}

	return b, nil
}

func IsAccountFollowing(ctx context.Context, db FollowingDatabase, account_id int64, following_address string) (bool, *Following, error) {

	f, err := db.GetFollowing(ctx, account_id, following_address)

	if err == nil {
		return true, f, nil
	}

	if err == ErrNotFound {
		return false, nil, nil
	}

	return false, nil, fmt.Errorf("Failed to following record, %w", err)
}
