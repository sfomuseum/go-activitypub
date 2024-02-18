package activitypub

type Note struct {
	Id            int64  `json:"id"`
	NoteId        string `json:"note_id"`
	AuthorAddress string `json:"author_address"`
	Body          []byte `json:"body"`
	Created       int64  `json:"created"`
	LastModified  int64  `json:"lastmodified"`
}
