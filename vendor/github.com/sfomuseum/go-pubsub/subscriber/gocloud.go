package subscriber

// https://gocloud.dev/howto/pubsub/subscribe/
// ./bin/subscribe -subscriber-uri 'awssqs-creds://?region={REGION}&credentials={CREDENTIALS}&queue-url=https://sqs.{REGION}.amazonaws.com/{ACCOUNT}/{QUEUE}'

import (
	"context"
	"fmt"
	"net/url"

	"github.com/aaronland/go-aws-auth"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"gocloud.dev/pubsub"
	"gocloud.dev/pubsub/awssnssqs"
)

type GoCloudSubscriber struct {
	Subscriber
	subscription *pubsub.Subscription
}

func init() {

	ctx := context.Background()
	RegisterGoCloudSubscribers(ctx)
}

func RegisterGoCloudSubscribers(ctx context.Context) error {

	to_register := []string{
		"awssqs-creds",
	}

	for _, scheme := range pubsub.DefaultURLMux().TopicSchemes() {
		to_register = append(to_register, scheme)
	}

	for _, scheme := range to_register {

		err := RegisterSubscriber(ctx, scheme, NewGoCloudSubscriber)

		if err != nil {
			return fmt.Errorf("Failed to register subscriber for '%s', %w", scheme, err)
		}
	}

	return nil
}

func NewGoCloudSubscriber(ctx context.Context, uri string) (Subscriber, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	var subscription *pubsub.Subscription

	switch u.Scheme {
	case "awssqs-creds":

		q := u.Query()

		region := q.Get("region")
		credentials := q.Get("credentials")
		queue_url := q.Get("queue-url")

		cfg_uri := fmt.Sprintf("aws://%s?credentials=%s", region, credentials)
		cfg, err := auth.NewConfig(ctx, cfg_uri)

		if err != nil {
			return nil, fmt.Errorf("Failed to create new session for credentials '%s', %w", credentials, err)
		}

		cl := sqs.NewFromConfig(cfg)

		if err != nil {
			return nil, fmt.Errorf("Failed to create AWS session, %w", err)
		}

		// https://gocloud.dev/howto/pubsub/publish/#sqs-ctor

		subscription = awssnssqs.OpenSubscriptionV2(ctx, cl, queue_url, nil)

	default:

		sub, err := pubsub.OpenSubscription(ctx, uri)

		if err != nil {
			return nil, err
		}

		subscription = sub
	}

	if err != nil {
		return nil, err
	}

	sub := &GoCloudSubscriber{
		subscription: subscription,
	}

	return sub, err
}

func (sub *GoCloudSubscriber) Listen(ctx context.Context, msg_ch chan string) error {

	for {

		msg, err := sub.subscription.Receive(ctx)

		if err != nil {
			return err
		}

		go msg.Ack()

		msg_ch <- string(msg.Body)
	}

	return nil
}

func (sub *GoCloudSubscriber) Close() error {
	ctx := context.Background()
	return sub.subscription.Shutdown(ctx)
}
