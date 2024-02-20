package activitypub

import (
	"strings"
)

const WEBFINGER_URI string = "/well-known/.webfinger"

type URIs struct {
	// Webfinger is assigned automatically
	Id       string `json:"id"`
	Activity string `json:"activity"`
	Profile  string `json:"profile"`
	Inbox    string `json:"inbox"`
	Outbox   string `json:"outbox"`
}

func DefaultURIs() *URIs {

	uris_table := &URIs{
		// Webfinger is assigned automatically
		Id:       "/ap/{resource}",
		Activity: "/ap/{resource}/activity",
		Profile:  "/ap/{resource}/profile",
		Inbox:    "/ap/{resource}/inbox",
		Outbox:   "/ap/{resource}/outbox",
	}

	return uris_table
}

func AssignResource(uri string, resource string) string {
	return strings.Replace(uri, "{resource}", resource, -1)
}
