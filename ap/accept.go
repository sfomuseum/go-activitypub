package ap

import (
	"context"

	"github.com/google/uuid"
)

func NewAcceptActivity(ctx context.Context, from string, to string) (*Activity, error) {

	guid := uuid.New()

	req := &Activity{
		Context: ACTIVITYSTREAMS_CONTEXT,
		Id:      guid.String(),
		Type:    "Accept",
		Actor:   from,
		Object:  to,
	}

	return req, nil
}
