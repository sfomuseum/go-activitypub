package activitypub

type Reply struct {
	Id        int64  `json:"id"`
	AccountId int64  `json:"account_id"`
	PostId    int64  `json:"post_id"`
	Actor     string `json:"actor"`
	ReplyId   string `json:"reply_id"`
	Body      string `json:"body"`
	Created   int64  `json:"created"`
}

func (r *Reply) Content() ([]byte, error) {
	return nil, ErrNotImplemented
}
