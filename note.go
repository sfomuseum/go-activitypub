package activitypub

type Note struct {
	Id           int64  `json:"id"`
	NoteId       string `json:"note_id"`
	AuthorId     string `json:"author_id"`
	Body         []byte `json:"body"`
	Created      int64  `json:"created"`
	LastModified int64  `json:"lastmodified"`
}
