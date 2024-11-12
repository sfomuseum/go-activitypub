package subscriber

import (
	"context"
	"fmt"
	"net/url"

	"github.com/redis/go-redis/v9"
	"github.com/sfomuseum/go-pubsub"
)

type RedisSubscriber struct {
	Subscriber
	redis_client  *redis.Client
	redis_channel string
}

func init() {
	ctx := context.Background()
	RegisterRedisSubscribers(ctx)
}

func RegisterRedisSubscribers(ctx context.Context) error {
	return RegisterSubscriber(ctx, "redis", NewRedisSubscriber)
}

func NewRedisSubscriber(ctx context.Context, uri string) (Subscriber, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	endpoint, channel, err := pubsub.RedisConfigFromURL(u)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive Redis config from URI, %w", err)
	}

	redis_client := redis.NewClient(&redis.Options{
		Addr: endpoint,
	})

	s := &RedisSubscriber{
		redis_client:  redis_client,
		redis_channel: channel,
	}

	return s, nil
}

func (s *RedisSubscriber) Listen(ctx context.Context, messages_ch chan string) error {

	pubsub_client := s.redis_client.PSubscribe(ctx, s.redis_channel)
	defer pubsub_client.Close()

	for {

		i, err := pubsub_client.Receive(ctx)

		if err != nil {
			return fmt.Errorf("Failed to receive message, %w", err)
		}

		if msg, _ := i.(*redis.Message); msg != nil {
			messages_ch <- msg.Payload
		}
	}

	return nil
}

func (s *RedisSubscriber) Close() error {
	return s.redis_client.Close()
}
