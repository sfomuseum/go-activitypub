package ap

import (
	"context"

	"github.com/sfomuseum/go-activitypub/id"
)

func NewCreateActivity(ctx context.Context, from string, to []string, object interface{}) (*Activity, error) {

	uuid := id.NewUUID()

	req := &Activity{
		Context: ACTIVITYSTREAMS_CONTEXT,
		Id:      uuid,
		Type:    "Create",
		Actor:   from,
		To:      to,
		Object:  object,
	}

	return req, nil

}
