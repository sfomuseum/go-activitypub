package activitypub

import (
	"github.com/sfomuseum/go-activitypub/ap"
)

const (
	// UndefinedObjectType is an unknown (or undefined) object type
	UndefinedObjectType ObjectType = iota
	PostObjectType
	BoostObjectType
)

type ObjectType int

// Type Activity is a (this) package-specific wrapper mapping ActivityPub activities to ...
type Activity struct {
	// A unique 64-bit ID for the activity
	Id int64 `json:"id"`
	// The unique ID associated with the ActivityPub Activity (Body)
	ActivityPubId string `json:"activity_id"`
	// The ...
	ActivityPubType string `json:"activity_type"`
	// The unique 64-bit ID for account associated with the activity
	AccountId int64 `json:"id"`
	// A valid object type (as in an ActivityPub object) supported by this package
	ObjectType ObjectType `json:"object_type"`
	// The unique 64-bit ID associated with the object type (for the activity)
	ObjectId int64 `json:"object_id"`
	// The body of the ActivityPub Activity
	Body *ap.Activity `json:"body"`
	// The Unix timestamp when the activity was created.
	Created int64 `json:"created"`
}
