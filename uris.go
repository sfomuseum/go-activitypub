package activitypub

type URIs struct {
	Profile  string `json:"profile"`
	Activity string `json:"activity"`
	Id       string `json:"id"`
	Inbox    string `json:"inbox"`
}

func DefaultURIs() *URIs {

	uris_table := &URIs{
		Profile:  "/profile/",
		Activity: "/actvity/",
		Id:       "/",
		Inbox:    "/inbox/",
	}

	return uris_table
}
