package queue

import (
	"context"
	"log/slog"
)

type SlogProcessFollowerQueue struct {
	ProcessFollowerQueue
}

func init() {
	ctx := context.Background()
	err := RegisterProcessFollowerQueue(ctx, "slog", NewSlogProcessFollowerQueue)

	if err != nil {
		panic(err)
	}
}

func NewSlogProcessFollowerQueue(ctx context.Context, uri string) (ProcessFollowerQueue, error) {
	q := &SlogProcessFollowerQueue{}
	return q, nil
}

func (q *SlogProcessFollowerQueue) ProcessFollower(ctx context.Context, follower_id int64) error {
	slog.Info("Process follower", "follower id", follower_id)
	return nil
}

func (q *SlogProcessFollowerQueue) Close(ctx context.Context) error {
	return nil
}
