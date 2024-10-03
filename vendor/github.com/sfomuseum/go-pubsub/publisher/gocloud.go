package publisher

// ./bin/publish -publisher-uri 'awssqs-creds://?region={REGION}&credentials={CREDENTIALS}&queue-url=https://sqs.{REGION}.amazonaws.com/{ACCOUNT}/{QUEUE}' 'hello world'

import (
	"context"
	"fmt"
	"net/url"

	"github.com/aaronland/go-aws-auth"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"gocloud.dev/pubsub"
	"gocloud.dev/pubsub/awssnssqs"
)

type GoCloudPublisher struct {
	Publisher
	topic *pubsub.Topic
}

func init() {

	ctx := context.Background()
	err := RegisterGoCloudPublishers(ctx)

	if err != nil {
		panic(err)
	}
}

// RegisterGoCloudPublishers will explicitly register all the schemes associated with the `GoCloudPublisher` interface.
func RegisterGoCloudPublishers(ctx context.Context) error {

	to_register := []string{
		"awssqs-creds",
	}

	for _, scheme := range pubsub.DefaultURLMux().TopicSchemes() {
		to_register = append(to_register, scheme)
	}

	for _, scheme := range to_register {

		err := RegisterPublisher(ctx, scheme, NewGoCloudPublisher)

		if err != nil {
			return fmt.Errorf("Failed to register blob publisher for '%s', %w", scheme, err)
		}
	}

	return nil
}

func NewGoCloudPublisher(ctx context.Context, uri string) (Publisher, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	var topic *pubsub.Topic

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

		topic = awssnssqs.OpenSQSTopicV2(ctx, cl, queue_url, nil)

	default:

		topic, err = pubsub.OpenTopic(ctx, uri)

		if err != nil {
			return nil, err
		}
	}

	pub := &GoCloudPublisher{
		topic: topic,
	}

	return pub, err
}

func (pub *GoCloudPublisher) Publish(ctx context.Context, str_msg string) error {

	msg := &pubsub.Message{
		Body: []byte(str_msg),
	}

	return pub.topic.Send(ctx, msg)
}

func (pub *GoCloudPublisher) Close() error {
	ctx := context.Background()
	return pub.topic.Shutdown(ctx)
}
