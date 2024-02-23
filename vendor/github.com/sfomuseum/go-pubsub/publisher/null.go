package publisher

import (
	"context"
)

type NullPublisher struct {
	Publisher
}

func init() {
	ctx := context.Background()
	RegisterNullPublishers(ctx)
}

func RegisterNullPublishers(ctx context.Context) error {
	return RegisterPublisher(ctx, "null", NewNullPublisher)
}

func NewNullPublisher(ctx context.Context, uri string) (Publisher, error) {

	pub := &NullPublisher{}
	return pub, nil
}

func (pub *NullPublisher) Publish(ctx context.Context, msg string) error {

	return nil
}

func (pub *NullPublisher) Close() error {
	return nil
}
