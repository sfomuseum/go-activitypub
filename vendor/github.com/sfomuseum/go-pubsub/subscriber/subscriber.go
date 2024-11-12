package subscriber

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"sync"

	"github.com/aaronland/go-roster"
)

// In principle this could also be done with a sync.OnceFunc call but that will
// require that everyone uses Go 1.21 (whose package import changes broke everything)
// which is literally days old as I write this. So maybe a few releases after 1.21.
//
// Also, _not_ using a sync.OnceFunc means we can call RegisterSchemes multiple times
// if and when multiple gomail-sender instances register themselves.

var register_mu = new(sync.RWMutex)
var register_map = map[string]bool{}

type Subscriber interface {
	Listen(context.Context, chan string) error
	Close() error
}

type SubscriberInitializeFunc func(ctx context.Context, uri string) (Subscriber, error)

var subscribers roster.Roster

func ensureSpatialRoster() error {

	if subscribers == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		subscribers = r
	}

	return nil
}

func RegisterSubscriber(ctx context.Context, scheme string, f SubscriberInitializeFunc) error {

	err := ensureSpatialRoster()

	if err != nil {
		return err
	}

	register_mu.Lock()
	defer register_mu.Unlock()

	_, exists := register_map[scheme]

	if exists {
		return nil
	}

	err = subscribers.Register(ctx, scheme, f)

	if err != nil {
		return err
	}

	register_map[scheme] = true
	return nil
}

func SubscriberSchemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureSpatialRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range subscribers.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}

func NewSubscriber(ctx context.Context, uri string) (Subscriber, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := subscribers.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	f := i.(SubscriberInitializeFunc)
	return f(ctx, uri)
}
