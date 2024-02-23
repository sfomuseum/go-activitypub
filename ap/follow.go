package ap

import (
	"context"
	"fmt"

	"github.com/sfomuseum/go-activitypub/id"
)

func NewFollowActivity(ctx context.Context, from string, to string) (*Activity, error) {

	uuid := id.NewUUID()

	req := &Activity{
		Id:     uuid,
		Type:   "Follow",
		Actor:  from,
		Object: to,
	}

	return req, nil
}

func NewUndoFollowActivity(ctx context.Context, from string, to string) (*Activity, error) {

	uuid := id.NewUUID()

	follow_activity, err := NewFollowActivity(ctx, from, to)

	if err != nil {
		return nil, fmt.Errorf("Failed to create follow activity (to undo), %w", err)
	}

	req := &Activity{
		Context: ACTIVITYSTREAMS_CONTEXT,
		Id:      uuid,
		Type:    "Undo",
		Actor:   from,
		Object:  follow_activity,
	}

	return req, nil
}
