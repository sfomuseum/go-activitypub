package subscriber

import (
	"context"
	"net/url"
	"time"
)

type TickerSubscriber struct {
	Subscriber
	ticker *time.Ticker
}

func init() {
	ctx := context.Background()
	RegisterTickerSubscribers(ctx)
}

func RegisterTickerSubscribers(ctx context.Context) error {
	return RegisterSubscriber(ctx, "ticker", NewTickerSubscriber)
}

func NewTickerSubscriber(ctx context.Context, uri string) (Subscriber, error) {

	_, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	t := time.NewTicker(1 * time.Second)

	sub := &TickerSubscriber{
		ticker: t,
	}

	return sub, nil
}

func (sub *TickerSubscriber) Listen(ctx context.Context, messages_ch chan string) error {

	for {
		select {
		case <-ctx.Done():
			return nil
		case t := <-sub.ticker.C:
			messages_ch <- t.Format(time.RFC3339)
		}
	}

	return nil
}

func (sub *TickerSubscriber) Close() error {
	sub.ticker.Stop()
	return nil
}
