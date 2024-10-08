package activitypub

import (
	"context"
	"fmt"
	"time"

	"github.com/sfomuseum/go-activitypub/id"
)

type Following struct {
	Id               int64  `json:"id"`
	AccountId        int64  `json:"account_id"`
	FollowingAddress string `json:"following_address"`
	Created          int64  `json:"created"`
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
