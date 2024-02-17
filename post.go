package activitypub

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sfomuseum/go-activitypub/ap"
)

type Post struct {
	Id           string `json:"id"`
	AccountId    string `json:"account_id"`
	Body         []byte `json:"body"`
	Created      int64  `json:"created"`
	LastModified int64  `json:"lastmodified"`
}

func NewPost(ctx context.Context, acct *Account, body []byte) (*Post, error) {

	now := time.Now()
	ts := now.Unix()

	guid := uuid.New()

	p := &Post{
		Id:           guid.String(),
		AccountId:    acct.Id,
		Body:         body,
		Created:      ts,
		LastModified: ts,
	}

	return p, nil
}

func (p *Post) AsNote(ctx context.Context) (*ap.Note, error) {

	// https://paul.kinlan.me/adding-activity-pub-to-your-static-site/

	t := time.Unix(p.Created, 0)

	n := &ap.Note{
		Type:         "Note",
		Id:           p.Id,
		AttributedTo: p.AccountId,
		To:           "https://www.w3.org/ns/activitystreams#Public", // what?
		Content:      string(p.Body),
		Published:    t.Format(time.RFC3339),
		URL:          "x-urn:fix-me",
	}

	return n, nil
}
