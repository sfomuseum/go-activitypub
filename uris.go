package activitypub

const WEBFINGER_URI string = "/well-known/.webfinger"

type URIs struct {
	// Webfinger is assigned automatically
	Profile  string `json:"profile"`
	Activity string `json:"activity"`
	Id       string `json:"id"`
	Inbox    string `json:"inbox"`
	Outbox   string `json:"outbox"`
}

func DefaultURIs() *URIs {

	uris_table := &URIs{
		// Webfinger is assigned automatically
		Profile:  "/profile/",
		Activity: "/actvity/",
		Id:       "/",
		Inbox:    "/inbox/",
		Outbox:   "/outbox/",
	}

	return uris_table
}
