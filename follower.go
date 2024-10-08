package activitypub

import (
	"context"
	"fmt"
	"time"

	"github.com/sfomuseum/go-activitypub/id"
)

type Follower struct {
	Id              int64  `json:"id"`
	AccountId       int64  `json:"account_id"`
	FollowerAddress string `json:"follower_address"`
	Created         int64  `json:"created"`
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
