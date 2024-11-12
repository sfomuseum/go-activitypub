package subscriber

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/hpcloud/tail" // Note: We are using a sfomuseum-specific fork by way of a replace directive in go.mod
)

type TailFileSubscriber struct {
	Subscriber
	path string
}

func init() {
	ctx := context.Background()
	RegisterTailFileSubscribers(ctx)
}

func RegisterTailFileSubscribers(ctx context.Context) error {
	return RegisterSubscriber(ctx, "tail", NewTailFileSubscriber)
}

func NewTailFileSubscriber(ctx context.Context, uri string) (Subscriber, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	sub := &TailFileSubscriber{
		path: u.Path,
	}

	return sub, nil
}

func (sub *TailFileSubscriber) Listen(ctx context.Context, messages_ch chan string) error {

	info, err := os.Stat(sub.path)

	if err != nil {
		return fmt.Errorf("Failed to stat %s, %w", sub.path, err)
	}

	seek_info := &tail.SeekInfo{
		Offset: info.Size(),
	}

	cfg := tail.Config{
		Location: seek_info,
		Follow:   true,
	}

	t, err := tail.TailFile(sub.path, cfg)

	if err != nil {
		return fmt.Errorf("Failed to tail %s, %w", sub.path, err)
	}

	for line := range t.Lines {
		messages_ch <- line.Text
	}

	err = t.Wait()

	if err != nil {
		return fmt.Errorf("Failed to wait tailing file, %w", err)
	}

	return nil
}

func (pub *TailFileSubscriber) Close() error {
	return nil
}
