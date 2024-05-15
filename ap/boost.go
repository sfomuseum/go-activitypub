package ap

import (
	"context"
	"fmt"
	"time"

	"github.com/sfomuseum/go-activitypub/uris"
)

// NewBoostActivity will return an ActivityPub "Announce" activity from 'from' about 'object' (created by 'author_uri').
func NewBoostActivity(ctx context.Context, uris_table *uris.URIs, from string, author_uri string, object interface{}) (*Activity, error) {
	return NewAnnounceActivity(ctx, uris_table, from, author_uri, object)
}

// NewAnnounceActivity will return an ActivityPub "Announce" activity from 'from' about 'object' (created by 'author_uri').
func NewAnnounceActivity(ctx context.Context, uris_table *uris.URIs, from string, author_uri string, object interface{}) (*Activity, error) {

	ap_id := NewId(uris_table)

	now := time.Now()

	activity := &Activity{
		Context: ACTIVITYSTREAMS_CONTEXT,
		Id:      ap_id,
		Type:    "Announce",
		Actor:   from,
		To: []string{
			fmt.Sprintf("%s#Public", ACTIVITYSTREAMS_CONTEXT),
		},
		Cc: []string{
			// Despite the example here which includes a URL it appears
			// that an address is necessary? See also: the DeliverPost
			// method in (sfomuseum/go-activitypub/deliver.go)
			// https://boyter.org/posts/activitypub-announce-post/
			author_uri,
		},
		Object:    object,
		Published: now.Format(time.RFC3339),
	}

	return activity, nil
}
