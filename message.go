package activitypub

type Message struct {
	Id           int64  `json:"id"`
	NoteId       int64  `json:"note_id"`
	AuthorURI    string `json:"author_uri"`
	AccountId    string `json:"account_id"`
	Created      int64  `json:"created"`
	LastModified int64  `json:"created"`
}
