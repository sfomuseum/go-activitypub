package queue

import (
	"context"
	"log/slog"
)

type SlogProcessFollowQueue struct {
	ProcessFollowQueue
}

func init() {
	ctx := context.Background()
	err := RegisterProcessFollowQueue(ctx, "slog", NewSlogProcessFollowQueue)

	if err != nil {
		panic(err)
	}
}

func NewSlogProcessFollowQueue(ctx context.Context, uri string) (ProcessFollowQueue, error) {
	q := &SlogProcessFollowQueue{}
	return q, nil
}

func (q *SlogProcessFollowQueue) ProcessFollow(ctx context.Context, account_id int64, follower string) error {
	slog.Info("Process follow", "account id", account_id, "follower", follower)
	return nil
}

func (q *SlogProcessFollowQueue) Close(ctx context.Context) error {
	return nil
}
