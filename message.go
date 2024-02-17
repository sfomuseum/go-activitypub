package activitypub

type Message struct {
	NoteId       int64  `json:"note_id"`
	AuthorId     string `json:"author_id"`
	AccountId    string `json:"account_id"`
	Created      int64  `json:"created"`
	LastModified int64  `json:"created"`
}
