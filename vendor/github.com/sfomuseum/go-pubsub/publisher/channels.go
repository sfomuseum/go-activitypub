package publisher

// https://eli.thegreenplace.net/2020/pubsub-using-channels-in-go/

import (
	"context"
	_ "log"
)

type ChannelPublisher struct {
	Publisher
	channel chan string
}

func NewChannelPublisherWithChannel(ctx context.Context, ch chan string) (Publisher, error) {

	pub := &ChannelPublisher{
		channel: ch,
	}

	return pub, nil
}

func (pub *ChannelPublisher) Publish(ctx context.Context, msg string) error {

	// Add timeout here...
	pub.channel <- msg
	return nil
}

func (pub *ChannelPublisher) Close() error {
	return nil
}
