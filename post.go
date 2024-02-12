package activitypub

type Post struct {
	Id   string `json:"id"`
	Body []byte `json:"body"`
}
