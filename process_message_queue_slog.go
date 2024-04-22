package activitypub

import (
	"context"
	"log/slog"
)

type SlogProcessMessageQueue struct {
	ProcessMessageQueue
}

func init() {
	ctx := context.Background()
	RegisterProcessMessageQueue(ctx, "slog", NewSlogProcessMessageQueue)
}

func NewSlogProcessMessageQueue(ctx context.Context, uri string) (ProcessMessageQueue, error) {
	q := &SlogProcessMessageQueue{}
	return q, nil
}

func (q *SlogProcessMessageQueue) ProcessMessage(ctx context.Context, message_id int64) error {
	slog.Info("Process message", "message_id", message_id)
	return nil
}
