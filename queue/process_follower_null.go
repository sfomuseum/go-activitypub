package queue

import (
	"context"
)

type NullProcessFollowerQueue struct {
	ProcessFollowerQueue
}

func init() {
	ctx := context.Background()
	err := RegisterProcessFollowerQueue(ctx, "null", NewNullProcessFollowerQueue)

	if err != nil {
		panic(err)
	}

}

func NewNullProcessFollowerQueue(ctx context.Context, uri string) (ProcessFollowerQueue, error) {
	q := &NullProcessFollowerQueue{}
	return q, nil
}

func (q *NullProcessFollowerQueue) ProcessFollower(ctx context.Context, follower_id int64) error {
	return nil
}

func (q *NullProcessFollowerQueue) Close(ctx context.Context) error {
	return nil
}
