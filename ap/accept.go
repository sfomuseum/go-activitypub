package ap

import (
	"context"

	"github.com/sfomuseum/go-activitypub/id"
)

func NewAcceptActivity(ctx context.Context, from string, to string) (*Activity, error) {

	uuid := id.NewUUID()

	req := &Activity{
		Context: ACTIVITYSTREAMS_CONTEXT,
		Id:      uuid,
		Type:    "Accept",
		Actor:   from,
		Object:  to,
	}

	return req, nil
}
