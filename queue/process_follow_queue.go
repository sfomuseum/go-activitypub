package queue

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/aaronland/go-roster"
)

type ProcessFollowQueue interface {
	ProcessFollow(context.Context, int64) error
	Close(context.Context) error
}

var process_follow_queue_roster roster.Roster

// ProcessFollowQueueInitializationFunc is a function defined by individual process_follow_queue package and used to create
// an instance of that process_follow_queue
type ProcessFollowQueueInitializationFunc func(ctx context.Context, uri string) (ProcessFollowQueue, error)

// RegisterProcessFollowQueue registers 'scheme' as a key pointing to 'init_func' in an internal lookup table
// used to create new `ProcessFollowQueue` instances by the `NewProcessFollowQueue` method.
func RegisterProcessFollowQueue(ctx context.Context, scheme string, init_func ProcessFollowQueueInitializationFunc) error {

	err := ensureProcessFollowQueueRoster()

	if err != nil {
		return err
	}

	return process_follow_queue_roster.Register(ctx, scheme, init_func)
}

func ensureProcessFollowQueueRoster() error {

	if process_follow_queue_roster == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		process_follow_queue_roster = r
	}

	return nil
}

// NewProcessFollowQueue returns a new `ProcessFollowQueue` instance configured by 'uri'. The value of 'uri' is parsed
// as a `url.URL` and its scheme is used as the key for a corresponding `ProcessFollowQueueInitializationFunc`
// function used to instantiate the new `ProcessFollowQueue`. It is assumed that the scheme (and initialization
// function) have been registered by the `RegisterProcessFollowQueue` method.
func NewProcessFollowQueue(ctx context.Context, uri string) (ProcessFollowQueue, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := process_follow_queue_roster.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	init_func := i.(ProcessFollowQueueInitializationFunc)
	return init_func(ctx, uri)
}

// Schemes returns the list of schemes that have been registered.
func ProcessFollowQueueSchemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureProcessFollowQueueRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range process_follow_queue_roster.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}
