package activitypub

type URIs struct {
	Webfinger string `json:"webfinger"`
	Profile   string `json:"profile"`
	Activity  string `json:"activity"`
	Id        string `json:"id"`
	Inbox     string `json:"inbox"`
	Outbox    string `json:"outbox"`
}

func DefaultURIs() *URIs {

	uris_table := &URIs{
		Webfinger: "/.webfinger",
		Profile:   "/profile/",
		Activity:  "/actvity/",
		Id:        "/",
		Inbox:     "/inbox/",
		Outbox:    "/outbox/",
	}

	return uris_table
}
