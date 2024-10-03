package ap

// https://www.w3.org/TR/activitystreams-vocabulary/

// Activity is a struct encapsulating an ActivityPub activity.
type Activity struct {
	// Context needs to be a "whatever" (interface{}) because ActivityPub (JSON-LD)
	// mixes and matches string URIs, arbritrary data structures and arrays of string
	// URIs and arbritrary data structures in @context...
	Context interface{} `json:"@context,omitempty"`
	// Id is the unique identifier for the activity.
	Id string `json:"id"`
	// Type is the name of the activity being performed.
	Type string `json:"type"`
	// Actor is the URI of the person (actor) performing the activity.
	Actor string `json:"actor"`
	// To is the list of URIs the activity should be delivered to.
	To []string `json:"to,omitempty"`
	// CC is the list of URIs the activity should be copied to.
	Cc []string `json:"cc,omitempty"`
	// Object is body of the activity itself.
	Object interface{} `json:"object,omitempty"`
}
