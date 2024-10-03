package ap

import (
	"context"

	"github.com/sfomuseum/go-activitypub/uris"
)

// NewCreateActivity returns a new `Activity` instance of type "Create".
func NewCreateActivity(ctx context.Context, uris_table *uris.URIs, from string, to []string, object interface{}) (*Activity, error) {

	ap_id := NewId(uris_table)

	req := &Activity{
		Context: ACTIVITYSTREAMS_CONTEXT,
		Id:      ap_id,
		Type:    CREATE_ACTIVITY,
		Actor:   from,
		To:      to,
		Object:  object,
	}

	return req, nil

}
