package ap

import (
	"context"

	"github.com/sfomuseum/go-activitypub/uris"
)

// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-accept

// NewAcceptActvity returns a new `Activity` instance of type "Accept".
// The Accept activity "indicates that the actor accepts the object. The target property can be used in certain circumstances to indicate the context into which the object has been accepted."
func NewAcceptActivity(ctx context.Context, uris_table *uris.URIs, from string, object interface{}) (*Activity, error) {

	ap_id := NewId(uris_table, "accept")

	req := &Activity{
		Context: []interface{}{
			ACTIVITYSTREAMS_CONTEXT,
		},
		Id:     ap_id,
		Type:   ACCEPT_ACTIVITY,
		Actor:  from,
		Object: object,
	}

	return req, nil
}
