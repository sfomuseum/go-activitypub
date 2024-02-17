package ap

import (
	"context"

	"github.com/google/uuid"
)

func NewCreateActivity(ctx context.Context, from string, to []string, object interface{}) (*Activity, error) {

	guid := uuid.New()

	req := &Activity{
		Context: ACTIVITYSTREAMS_CONTEXT,
		Id:      guid.String(),
		Type:    "Create",
		Actor:   from,
		To:      to,
		Object:  object,
	}

	return req, nil

}
