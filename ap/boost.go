package ap

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/sfomuseum/go-activitypub/uris"
)

// TBD replace the return value ap.Activity to be activitypub.Activity ?
// Also move this in to go-activitypub/boosts...

func NewBoostActivityForNote(ctx context.Context, uris_table *uris.URIs, from string, note_uri string) (*Activity, error) {

	logger := slog.Default()
	logger = logger.With("from", from)
	logger = logger.With("note", note_uri)

	logger.Debug("Retrieve note")

	n, err := RetrieveNote(ctx, note_uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve note, %w", err)
	}

	logger = logger.With("attributed to", n.AttributedTo)
	logger.Debug("Retrieve note")

	author, err := RetrieveActorWithProfileURL(ctx, n.AttributedTo)

	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve author (actor) for note, %w", err)
	}

	logger.Debug("Derive address for author")

	author_addr, err := author.Address()

	if err != nil {
		return nil, fmt.Errorf("Failed to derive note author address, %w", err)
	}

	logger = logger.With("author address", author_addr)
	logger.Debug("Create new boost (annouce) activity")

	return NewBoostActivity(ctx, uris_table, from, author_addr, note_uri)
}

// NewBoostActivity will return an ActivityPub "Announce" activity from 'from' about 'object' (created by 'author_addr').
func NewBoostActivity(ctx context.Context, uris_table *uris.URIs, from string, author_addr string, object interface{}) (*Activity, error) {
	return NewAnnounceActivity(ctx, uris_table, from, author_addr, object)
}

// NewAnnounceActivity will return an ActivityPub "Announce" activity from 'from' about 'object' (created by 'author_addr').
func NewAnnounceActivity(ctx context.Context, uris_table *uris.URIs, from string, author_addr string, object interface{}) (*Activity, error) {

	ap_id := NewId(uris_table, "announce")

	now := time.Now()

	activity := &Activity{
		Context: ACTIVITYSTREAMS_CONTEXT,
		Id:      ap_id,
		Type:    "Announce",
		Actor:   from,
		To: []string{
			ACTIVITYSTREAMS_CONTEXT_PUBLIC,
		},
		Cc: []string{
			// Despite the example here which includes a URL it appears
			// that an address is necessary? See also: the DeliverActivity
			// method in (sfomuseum/go-activitypub/queue/deliver_activity.go)
			// https://boyter.org/posts/activitypub-announce-post/
			author_addr,
		},
		Object:    object,
		Published: now.Format(time.RFC3339),
	}

	return activity, nil
}
