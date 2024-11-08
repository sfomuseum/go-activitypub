package queue

import (
	"context"
)

type NullProcessMessageQueue struct {
	ProcessMessageQueue
}

func init() {
	ctx := context.Background()
	err := RegisterProcessMessageQueue(ctx, "null", NewNullProcessMessageQueue)

	if err != nil {
		panic(err)
	}

}

func NewNullProcessMessageQueue(ctx context.Context, uri string) (ProcessMessageQueue, error) {
	q := &NullProcessMessageQueue{}
	return q, nil
}

func (q *NullProcessMessageQueue) ProcessMessage(ctx context.Context, message_id int64) error {
	return nil
}

func (q *NullProcessMessageQueue) Close(ctx context.Context) error {
	return nil
}
