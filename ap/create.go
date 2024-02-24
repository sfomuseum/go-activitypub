package ap

import (
	"context"

	"github.com/sfomuseum/go-activitypub/uris"
)

func NewCreateActivity(ctx context.Context, uris_table *uris.URIs, from string, to []string, object interface{}) (*Activity, error) {

	ap_id := NewId(uris_table)

	req := &Activity{
		Context: ACTIVITYSTREAMS_CONTEXT,
		Id:      ap_id,
		Type:    "Create",
		Actor:   from,
		To:      to,
		Object:  object,
	}

	return req, nil

}
