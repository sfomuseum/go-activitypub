package activitypub

import (
	"net/url"
	"strings"
)

const WEBFINGER_URI string = "/well-known/.webfinger"

type URIs struct {
	// Webfinger is assigned automatically

	Account   string `json:"account"`
	Inbox     string `json:"inbox"`
	Outbox    string `json:"outbox"`
	Followers string `json:"followers"`
	Following string `json:"following"`
	Icon      string `json:"icon"`

	Hostname string `json:"hostname"`
	Insecure bool   `json:"insecure"`
}

func DefaultURIs() *URIs {

	uris_table := &URIs{
		// Webfinger is assigned automatically

		Account:   "/ap/{resource}",
		Inbox:     "/ap/{resource}/inbox",
		Outbox:    "/ap/{resource}/outbox",
		Followers: "/ap/{resource}/followers",
		Following: "/ap/{resource}/following",
		Icon:      "/ap/{resource}/icon.png",
	}

	return uris_table
}

func AssignResource(uri string, resource string) string {
	return strings.Replace(uri, "{resource}", resource, -1)
}

func NewURL(uris_table *URIs, path string) *url.URL {

	scheme := "https"

	if uris_table.Insecure {
		scheme = "http"
	}

	host := uris_table.Hostname

	u := &url.URL{}
	u.Scheme = scheme
	u.Host = host
	u.Path = path

	return u
}
