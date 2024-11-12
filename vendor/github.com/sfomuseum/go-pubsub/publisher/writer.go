package publisher

import (
	"context"
	"io"
)

// PublisherWriter wraps a `Publisher` instance and confirms to the `io.WriteCloser` interface.
type PublisherWriter struct {
	io.WriteCloser
	publisher Publisher
}

// NewWriter returns a new instance of `PublisherWriter` which wraps 'p' and confirms to the `io.WriteCloser` interface.
func NewWriter(p Publisher) io.WriteCloser {

	pwr := &PublisherWriter{
		publisher: p,
	}

	return pwr
}

// Write dispatches 'p' to the underlyng `Publisher` instance.
func (pwr *PublisherWriter) Write(p []byte) (int, error) {

	ctx := context.Background()
	err := pwr.publisher.Publish(ctx, string(p))

	if err != nil {
		return 0, err
	}

	return len(p), nil
}

// Closer triggers the underlyng `Publisher` instance's Close method.
func (pwr *PublisherWriter) Close() error {
	return pwr.publisher.Close()
}
