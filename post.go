package activitypub

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type Post struct {
	Id           string `json:"id"`
	AccountId    string `json:"account_id"`
	Body         []byte `json:"body"`
	Created      int64  `json:"created"`
	LastModified int64  `json:"lastmodified"`
}

func NewPost(ctx context.Context, acct *Account, body []byte) (*Post, error) {

	now := time.Now()
	ts := now.Unix()

	guid := uuid.New()

	p := &Post{
		Id:           guid.String(),
		AccountId:    acct.Id,
		Body:         body,
		Created:      ts,
		LastModified: ts,
	}

	return p, nil
}

func (p *Post) Deliver(ctx context.Context, followers_db FollowersDatabase, q DeliveryQueue) error {

	followers_cb := func(ctx context.Context, follower_id string) error {

		slog.Info("Deliver", "post", p.Id, "follower_id", follower_id)
		err := q.DeliverPost(ctx, p, follower_id)

		if err != nil {
			return fmt.Errorf("Failed to deliver post to %s, %w", follower_id, err)
		}

		return nil
	}

	err := followers_db.GetFollowers(ctx, p.AccountId, followers_cb)

	if err != nil {
		return fmt.Errorf("Failed to get followers for post author, %w", err)
	}

	return nil
}
