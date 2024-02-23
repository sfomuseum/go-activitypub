// package publisher provides a common interface for publish operations.
package publisher

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

type Publisher interface {
	Publish(context.Context, string) error
	Close() error
}

type PublisherInitializeFunc func(ctx context.Context, uri string) (Publisher, error)

var publishers roster.Roster

func ensurePublisherRoster() error {

	if publishers == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		publishers = r
	}

	return nil
}

func RegisterPublisher(ctx context.Context, scheme string, f PublisherInitializeFunc) error {

	err := ensurePublisherRoster()

	if err != nil {
		return err
	}

	register_mu.Lock()
	defer register_mu.Unlock()

	_, exists := register_map[scheme]

	if exists {
		return nil
	}

	err = publishers.Register(ctx, scheme, f)

	if err != nil {
		return err
	}

	register_map[scheme] = true
	return nil
}

func PublisherSchemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensurePublisherRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range publishers.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}

func NewPublisher(ctx context.Context, uri string) (Publisher, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := publishers.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	f := i.(PublisherInitializeFunc)
	return f(ctx, uri)
}
