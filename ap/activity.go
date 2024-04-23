package ap

// https://www.w3.org/TR/activitystreams-vocabulary/

type Activity struct {
	// Context needs to be a "whatever" (interface{}) because ActivityPub (JSON-LD)
	// mixes and matches string URIs, arbritrary data structures and arrays of string
	// URIs and arbritrary data structures in @context...
	Context interface{} `json:"@context,omitempty"`
	Id      string      `json:"id"`
	Type    string      `json:"type"`
	Actor   string      `json:"actor"`
	To      []string    `json:"to,omitempty"`
	Cc      []string    `json:"cc,omitempty"`
	Object  interface{} `json:"object,omitempty"`
}
