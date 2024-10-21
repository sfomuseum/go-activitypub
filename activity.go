package activitypub

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sfomuseum/go-activitypub/ap"
	"github.com/sfomuseum/go-activitypub/id"
)

const (
	// UndefinedActivityType is an unknown (or undefined) object type
	UndefinedActivityType ActivityType = iota
	PostActivityType
	BoostActivityType
)

type ActivityType int

// Type Activity is internal representation of an ActivityPub Activity message with pointers
// to other relevant internal representations of things like account IDs rather than actor
// URIs, etc.
type Activity struct {
	// A unique 64-bit ID for the activity
	Id int64 `json:"id"`
	// The unique ID associated with the ActivityPub Activity (Body)
	ActivityPubId string `json:"activity_id"`
	// The unique 64-bit ID for account associated with the activity
	AccountId int64 `json:"id"`
	// A valid object type (as in an ActivityPub object) supported by this package
	ActivityType ActivityType `json:"activity_type"`
	// The unique 64-bit ID associated with the activity type (for example if ActivityType is PostActivityType then ActivityId would be unique ID for that post).
	ActivityTypeId int64 `json:"activity_type_id"`
	// The JSON-encoded body of the ActivityPub Activity
	Body string `json:"body"`
	// The Unix timestamp when the activity was created.
	Created int64 `json:"created"`
}

// NewActivity returns a new `Activity` instance using properties derived from 'ap_activity'.
func NewActivity(ctx context.Context, ap_activity *ap.Activity) (*Activity, error) {

	enc_ap, err := json.Marshal(ap_activity)

	if err != nil {
		return nil, fmt.Errorf("Failed to marshal activity, %w", err)
	}

	id, err := id.NewId()

	if err != nil {
		return nil, fmt.Errorf("Failed to create new ID, %w", err)
	}

	ap_id := ap_activity.Id

	now := time.Now()
	ts := now.Unix()

	a := &Activity{
		Id:            id,
		ActivityPubId: ap_id,
		Body:          string(enc_ap),
		Created:       ts,
	}

	return a, nil
}

// UnmarshalActivity returns the unmarshaled `*ap.Activity` that is encapsulated by 'a'.
func (a *Activity) UnmarshalActivity() (*ap.Activity, error) {

	var ap_activity *ap.Activity
	err := json.Unmarshal([]byte(a.Body), &ap_activity)

	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal activity, %w", err)
	}

	return ap_activity, nil
}
