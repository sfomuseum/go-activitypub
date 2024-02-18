package activitypub

type Message struct {
	Id            int64  `json:"id"`
	NoteId        int64  `json:"note_id"`
	AuthorAddress string `json:"author_uri"`
	AccountId     int64  `json:"account_id"`
	Created       int64  `json:"created"`
	LastModified  int64  `json:"created"`
}
