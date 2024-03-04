package activitypub

import (
	"context"
	"fmt"
	"time"

	"github.com/sfomuseum/go-activitypub/id"
)

type Boost struct {
	Id        int64  `json:"id"`
	AccountId int64  `json:"account_id"`
	PostId    int64  `json:"post_id"`
	Actor     string `json:"actor"`
	Created   int64  `json:"created"`
}

func NewBoost(ctx context.Context, post *Post, actor string) (*Boost, error) {

	boost_id, err := id.NewId()

	if err != nil {
		return nil, fmt.Errorf("Failed to create new ID, %w", err)
	}

	now := time.Now()
	ts := now.Unix()

	l := &Boost{
		Id:        boost_id,
		AccountId: post.AccountId,
		PostId:    post.Id,
		Actor:     actor,
		Created:   ts,
	}

	return l, nil
}
