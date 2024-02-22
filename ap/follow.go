package ap

import (
	"context"

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

func NewUnFollowActivity(ctx context.Context, from string, to string) (*Activity, error) {

	uuid := id.NewUUID()

	req := &Activity{
		Context: ACTIVITYSTREAMS_CONTEXT,
		Id:      uuid,
		Type:    "Undo",
		Actor:   from,
		Object:  to,
	}

	return req, nil
}
