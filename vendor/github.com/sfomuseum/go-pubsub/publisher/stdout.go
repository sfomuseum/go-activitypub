package publisher

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strconv"
)

type StdoutPublisher struct {
	Publisher
	apply_newline bool
}

func init() {
	ctx := context.Background()
	RegisterStdoutPublishers(ctx)
}

func RegisterStdoutPublishers(ctx context.Context) error {
	return RegisterPublisher(ctx, "stdout", NewStdoutPublisher)
}

func NewStdoutPublisher(ctx context.Context, uri string) (Publisher, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	q := u.Query()

	pub := &StdoutPublisher{}

	if q.Get("newline") != "" {

		v, err := strconv.ParseBool(q.Get("newline"))

		if err != nil {
			return nil, fmt.Errorf("Failed to parse ?newline= parameter, %w", err)
		}

		pub.apply_newline = v
	}

	return pub, nil
}

func (pub *StdoutPublisher) Publish(ctx context.Context, msg string) error {

	if pub.apply_newline {
		msg = fmt.Sprintf("%s\n", msg)
	}

	os.Stdout.Write([]byte(msg))
	return nil
}

func (pub *StdoutPublisher) Close() error {
	return nil
}
