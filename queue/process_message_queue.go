package queue

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/aaronland/go-roster"
)

type ProcessMessageQueue interface {
	ProcessMessage(context.Context, int64) error
	Close(context.Context) error
}

var process_message_queue_roster roster.Roster

// ProcessMessageQueueInitializationFunc is a function defined by individual process_message_queue package and used to create
// an instance of that process_message_queue
type ProcessMessageQueueInitializationFunc func(ctx context.Context, uri string) (ProcessMessageQueue, error)

// RegisterProcessMessageQueue registers 'scheme' as a key pointing to 'init_func' in an internal lookup table
// used to create new `ProcessMessageQueue` instances by the `NewProcessMessageQueue` method.
func RegisterProcessMessageQueue(ctx context.Context, scheme string, init_func ProcessMessageQueueInitializationFunc) error {

	err := ensureProcessMessageQueueRoster()

	if err != nil {
		return err
	}

	return process_message_queue_roster.Register(ctx, scheme, init_func)
}

func ensureProcessMessageQueueRoster() error {

	if process_message_queue_roster == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		process_message_queue_roster = r
	}

	return nil
}

// NewProcessMessageQueue returns a new `ProcessMessageQueue` instance configured by 'uri'. The value of 'uri' is parsed
// as a `url.URL` and its scheme is used as the key for a corresponding `ProcessMessageQueueInitializationFunc`
// function used to instantiate the new `ProcessMessageQueue`. It is assumed that the scheme (and initialization
// function) have been registered by the `RegisterProcessMessageQueue` method.
func NewProcessMessageQueue(ctx context.Context, uri string) (ProcessMessageQueue, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := process_message_queue_roster.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	init_func := i.(ProcessMessageQueueInitializationFunc)
	return init_func(ctx, uri)
}

// Schemes returns the list of schemes that have been registered.
func ProcessMessageQueueSchemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureProcessMessageQueueRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range process_message_queue_roster.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}
