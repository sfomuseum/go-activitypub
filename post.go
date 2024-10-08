package activitypub

import (
	"context"
	"fmt"
	"time"

	"github.com/sfomuseum/go-activitypub/id"
)

// Type Post is probably a misnomer. Specifically it began life as the internal representation
// of a "post" or a generic "message" distinct from the ActvityPub "activity" type stored in a
// "posts" database (see PostsDatabase). To date all the other code – notably the code to deliver
// AP messages to other servers – has been pretty tightly wrapped in this idea. For example in
// deliver.go we are explictly calling NoteFromPost(ctx, opts.URIs, opts.From, opts.Post, opts.PostTags)
// and then ap.NewCreateActivity(ctx, opts.URIs, from_uri, to_list, note).
//
// So basically the options are:
//
//  1. Treat the "Body" text as "special" and trap and interpret specific kinds of messages (in places
//     like deliver.go)
//  2. Update the Post struct (and databases) to include something about "boosts" (and eventually other
//     types of messages, for example "likes"... at which point you start to better understand the notion
//     of "activities"...)
//  3. Update the Boost struct (and databases) to denote directionality. Then do the same thing for the likes.
//  4. Update deliver.go to rename DeliverPost options and methods to be DeliverMessage and add "Boost" and
//     "Like" properties and use those to create relevant AP activities (in deliver.go)
//
// None of these are great. As I write this (20240430) I am inclined to favour (1) for the following reasons:
//
//  1. It means that we still have a record of all the actions (activities) performed by an account unlike
//     option (4) for example
//  2. Assuming a well-defined and consistent syntax for denoting post bodies that are not notes/create activities
//     then there is a path for deconstructing or updating the Post struct (and PostsDatabase schema) at a later
//     date whether that involves creating a new Actions/Activities database or updating the Likes/Boosts databases.
//  3. All of the special-case logic is confined to deliver.go and strictly-enforced conventions for doing boosts
//     or likes, for example boost:{ACCOUNT_BEING_BOOSTED}:{URI_BEING_BOOSTED} or boost://?{PARAMS} and so on
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
