package activitypub

import (
	"context"

	"github.com/sfomuseum/go-activitypub/ap"
)

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

func NewReply(ctx context.Context, note *ap.Note, post *Post) (*Reply, error) {
	return nil, ErrNotImplemented
}
