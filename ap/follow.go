package ap

import (
	"context"
	"fmt"

	"github.com/sfomuseum/go-activitypub/uris"
)

func NewFollowActivity(ctx context.Context, uris_table *uris.URIs, from string, to string) (*Activity, error) {

	ap_id := NewId(uris_table)

	req := &Activity{
		Id:     ap_id,
		Type:   "Follow",
		Actor:  from,
		Object: to,
	}

	return req, nil
}

func NewUndoFollowActivity(ctx context.Context, uris_table *uris.URIs, from string, to string) (*Activity, error) {

	ap_id := NewId(uris_table)

	follow_activity, err := NewFollowActivity(ctx, uris_table, from, to)

	if err != nil {
		return nil, fmt.Errorf("Failed to create follow activity (to undo), %w", err)
	}

	req := &Activity{
		Context: ACTIVITYSTREAMS_CONTEXT,
		Id:      ap_id,
		Type:    "Undo",
		Actor:   from,
		Object:  follow_activity,
	}

	return req, nil
}
