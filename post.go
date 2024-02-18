package activitypub

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sfomuseum/go-activitypub/ap"
)

type Post struct {
	Id           int64  `json:"id"`
	UUID         string `json:"uuid"`
	AccountId    int64  `json:"account_id"`
	Body         []byte `json:"body"`
	Created      int64  `json:"created"`
	LastModified int64  `json:"lastmodified"`
}

func NewPost(ctx context.Context, acct *Account, body []byte) (*Post, error) {

	post_id, err := NewId()

	if err != nil {
		return nil, fmt.Errorf("Failed to derive new post ID, %w", err)
	}

	now := time.Now()
	ts := now.Unix()

	guid := uuid.New()

	p := &Post{
		Id:           post_id,
		UUID:         guid.String(),
		AccountId:    acct.Id,
		Body:         body,
		Created:      ts,
		LastModified: ts,
	}

	return p, nil
}

func (p *Post) AsNote(ctx context.Context) (*ap.Note, error) {

	// https://paul.kinlan.me/adding-activity-pub-to-your-static-site/

	// Need hostname and URIs and possible accounts database?
	url := fmt.Sprintf("x-urn:fix-me#%d", p.Id)

	// Need account or accounts database...
	attr := "fix me"

	guid := p.UUID

	t := time.Unix(p.Created, 0)

	n := &ap.Note{
		Type:         "Note",
		Id:           guid,
		AttributedTo: attr,
		To:           "https://www.w3.org/ns/activitystreams#Public", // what?
		Content:      string(p.Body),
		Published:    t.Format(time.RFC3339),
		URL:          url,
	}

	return n, nil
}
