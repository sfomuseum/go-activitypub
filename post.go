package activitypub

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/sfomuseum/go-activitypub/ap"
	"github.com/sfomuseum/go-activitypub/id"
	"github.com/sfomuseum/go-activitypub/uris"
)

type Post struct {
	Id        int64  `json:"id"`
	UUID      string `json:"uuid"`
	AccountId int64  `json:"account_id"`
	// This is a string mostly because []byte thingies get encoded incorrectly
	// in DynamoDB
	Body         string `json:"body"`
	Created      int64  `json:"created"`
	LastModified int64  `json:"lastmodified"`
}

func NewPost(ctx context.Context, acct *Account, body string) (*Post, error) {

	post_id, err := id.NewId()

	if err != nil {
		return nil, fmt.Errorf("Failed to derive new post ID, %w", err)
	}

	uuid := id.NewUUID()

	now := time.Now()
	ts := now.Unix()

	p := &Post{
		Id:           post_id,
		UUID:         uuid,
		AccountId:    acct.Id,
		Body:         body,
		Created:      ts,
		LastModified: ts,
	}

	return p, nil
}

func NoteFromPost(ctx context.Context, uris_table *uris.URIs, acct *Account, post *Post) (*ap.Note, error) {

	// Need account or accounts database...
	attr := acct.ProfileURL(ctx, uris_table).String()

	post_url := acct.PostURL(ctx, uris_table, post)

	ap_id := ap.NewId(uris_table)

	t := time.Unix(post.Created, 0)

	n := &ap.Note{
		Type:         "Note",
		Id:           ap_id,
		AttributedTo: attr,
		To:           "https://www.w3.org/ns/activitystreams#Public", // what?
		Content:      post.Body,
		Published:    t.Format(http.TimeFormat),
		URL:          post_url.String(),
	}

	return n, nil
}
