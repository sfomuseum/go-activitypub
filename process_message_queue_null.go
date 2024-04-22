package activitypub

import (
	"context"
)

type NullProcessMessageQueue struct {
	ProcessMessageQueue
}

func init() {
	ctx := context.Background()
	RegisterProcessMessageQueue(ctx, "null", NewNullProcessMessageQueue)
}

func NewNullProcessMessageQueue(ctx context.Context, uri string) (ProcessMessageQueue, error) {
	q := &NullProcessMessageQueue{}
	return q, nil
}

func (q *NullProcessMessageQueue) ProcessMessage(ctx context.Context, message_id int64) error {
	return nil
}
