package ap

import (
	"context"

	"github.com/sfomuseum/go-activitypub/uris"
)

func NewAcceptActivity(ctx context.Context, uris_table *uris.URIs, from string, object interface{}) (*Activity, error) {

	ap_id := NewId(uris_table)

	req := &Activity{
		Context: []interface{}{
			ACTIVITYSTREAMS_CONTEXT,
		},
		Id:     ap_id,
		Type:   "Accept",
		Actor:  from,
		Object: object,
	}

	return req, nil
}
