package subscriber

import (
	"context"
)

type ChannelSubscriber struct {
	Subscriber
	channel chan string
}

func NewChannelSubscriberWithChannel(ctx context.Context, ch chan string) (Subscriber, error) {

	sub := &ChannelSubscriber{
		channel: ch,
	}

	return sub, nil
}

func (sub *ChannelSubscriber) Listen(ctx context.Context, messages_ch chan string) error {

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg := <-sub.channel:
			messages_ch <- msg
		}
	}

	return nil
}

func (sub *ChannelSubscriber) Close() error {
	return nil
}
