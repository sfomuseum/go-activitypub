package activitypub

import (
	"context"
	"fmt"
	"time"

	"github.com/sfomuseum/go-activitypub/id"
)

// Post is a message (or post) written by an account holder. It is internal representation
// of what would be delivered as an ActivityPub note.x
type Post struct {
	// The unique ID for the post.
	Id int64 `json:"id"`
	// The AccountsDatabase ID of the author of the post.
	AccountId int64 `json:"account_id"`
	// The body of the post. This is a string mostly because []byte thingies get encoded incorrectly
	// in DynamoDB
	Body string `json:"body"`
	// The URL of the post this post is referencing.
	InReplyTo string `json:"in_reply_to"`
	// The Unix timestamp when the post was created
	Created int64 `json:"created"`
	// The Unix timestamp when the post was last modified
	LastModified int64 `json:"lastmodified"`
}

// NewPost returns a new `Post` instance from 'acct' and 'body'.
func NewPost(ctx context.Context, acct *Account, body string) (*Post, error) {

	post_id, err := id.NewId()

	if err != nil {
		return nil, fmt.Errorf("Failed to derive new post ID, %w", err)
	}

	now := time.Now()
	ts := now.Unix()

	p := &Post{
		Id:           post_id,
		AccountId:    acct.Id,
		Body:         body,
		Created:      ts,
		LastModified: ts,
	}

	return p, nil
}
