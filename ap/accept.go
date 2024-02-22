package ap

import (
	"context"

	"github.com/sfomuseum/go-activitypub/id"
)

func NewAcceptActivity(ctx context.Context, from string, object interface{}) (*Activity, error) {

	uuid := id.NewUUID()

	req := &Activity{
		Context: ACTIVITYSTREAMS_CONTEXT,
		Id:      uuid,
		Type:    "Accept",
		Actor:   from,
		Object:  object,
	}

	return req, nil
}
