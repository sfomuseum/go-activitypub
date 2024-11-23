package queue

import (
	"context"
)

type NullProcessFollowQueue struct {
	ProcessFollowQueue
}

func init() {
	ctx := context.Background()
	err := RegisterProcessFollowQueue(ctx, "null", NewNullProcessFollowQueue)

	if err != nil {
		panic(err)
	}

}

func NewNullProcessFollowQueue(ctx context.Context, uri string) (ProcessFollowQueue, error) {
	q := &NullProcessFollowQueue{}
	return q, nil
}

func (q *NullProcessFollowQueue) ProcessFollow(ctx context.Context, follower_id int64) error {
	return nil
}

func (q *NullProcessFollowQueue) Close(ctx context.Context) error {
	return nil
}
