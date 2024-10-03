package publisher

import (
	"context"
	"fmt"
	"net/url"

	"github.com/redis/go-redis/v9"
	"github.com/sfomuseum/go-pubsub"
)

type RedisPublisher struct {
	Publisher
	redis_client  *redis.Client
	redis_channel string
}

func init() {
	ctx := context.Background()
	RegisterRedisPublishers(ctx)
}

func RegisterRedisPublishers(ctx context.Context) error {
	return RegisterPublisher(ctx, "redis", NewRedisPublisher)
}

func NewRedisPublisher(ctx context.Context, uri string) (Publisher, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	endpoint, channel, err := pubsub.RedisConfigFromURL(u)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive Redis config from URI, %w", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: endpoint,
	})

	p := &RedisPublisher{
		redis_client:  client,
		redis_channel: channel,
	}

	return p, nil
}

func (p *RedisPublisher) Publish(ctx context.Context, msg string) error {

	select {
	case <-ctx.Done():
		return nil
	default:
		// pass
	}

	cmd := p.redis_client.Publish(ctx, p.redis_channel, msg)
	err := cmd.Err()

	if err != nil {
		return fmt.Errorf("Failed to publish message, %w", err)
	}

	return nil
}

func (p *RedisPublisher) Close() error {
	return p.redis_client.Close()
}
