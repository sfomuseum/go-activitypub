package activitypub

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/sfomuseum/go-activitypub/ap"
	"github.com/sfomuseum/go-activitypub/id"
	"github.com/sfomuseum/go-activitypub/uris"
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

type AddPostOptions struct {
	URIs             *uris.URIs
	PostsDatabase    PostsDatabase
	PostTagsDatabase PostTagsDatabase
}

// AddPost creates a new post record for 'body' and adds it to the post database. Then it parses 'body' looking
// for other ActivityPub addresses and records each as a "mention" in the post tags database. It returns the post
// and the list of post tags (mentions) for further processing as needed.
func AddPost(ctx context.Context, opts *AddPostOptions, acct *Account, body string) (*Post, []*PostTag, error) {

	// Create the new post record

	p, err := NewPost(ctx, acct, body)

	if err != nil {
		return nil, nil, fmt.Errorf("Failed to create new post, %w", err)
	}

	// Add the post to the database

	err = opts.PostsDatabase.AddPost(ctx, p)

	if err != nil {
		return nil, nil, fmt.Errorf("Failed to add post, %w", err)
	}

	// Determine other accounts mentioned in post

	addrs_mentioned, err := ParseAddressesFromString(body)

	if err != nil {
		return nil, nil, fmt.Errorf("Failed to derive addresses mentioned in message body, %w", err)
	}

	// Create "mention" tags for any other accounts mentioned in post
	// Add each mention to the "post tags" database

	post_tags := make([]*PostTag, 0)

	for _, name := range addrs_mentioned {

		actor, err := RetrieveActor(ctx, name, opts.URIs.Insecure)

		if err != nil {
			slog.Error("Failed to retrieve actor data for name, skipping", "name", name, "error", err)
			continue
		}

		mention_name := name

		// https://github.com/sfomuseum/go-activitypub/issues/3
		// mention_href := actor.URL

		// And yet it appears to actually be {ACTOR}.id however this
		// does not work (where "work" means open profile tab) in Ivory
		// yet because... I have no idea
		mention_href := actor.Id

		t, err := NewMention(ctx, p, mention_name, mention_href)

		if err != nil {
			return nil, nil, fmt.Errorf("Failed to create mention for '%s', %w", name, err)
		}

		err = opts.PostTagsDatabase.AddPostTag(ctx, t)

		if err != nil {
			return nil, nil, fmt.Errorf("Failed to record post tag (mention) for '%s', %w", name, err)
		}

		post_tags = append(post_tags, t)
	}

	// Return all the things

	return p, post_tags, nil
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

// NoteFromPost creates a new (ActivityPub) `Note` instance derived from 'acct', 'post' and 'post_tags'.
func NoteFromPost(ctx context.Context, uris_table *uris.URIs, acct *Account, post *Post, post_tags []*PostTag) (*ap.Note, error) {

	attr := acct.ProfileURL(ctx, uris_table).String()
	post_url := acct.PostURL(ctx, uris_table, post)

	t := time.Unix(post.Created, 0)

	tags := make([]*ap.Tag, len(post_tags))
	cc := make([]string, 0)

	for idx, pt := range post_tags {

		t := &ap.Tag{
			Name: pt.Name,
			Href: pt.Href,
			Type: pt.Type,
		}

		tags[idx] = t

		if pt.Type == "Mention" {
			cc = append(cc, pt.Href)
		}
	}

	n := &ap.Note{
		Type:         "Note",
		Id:           post_url.String(),
		AttributedTo: attr,
		To: []string{
			"https://www.w3.org/ns/activitystreams#Public", // what?
		},
		Content:   post.Body,
		Published: t.Format(http.TimeFormat),
		InReplyTo: post.InReplyTo,
		Tags:      tags,
		URL:       post_url.String(),
	}

	if len(cc) > 0 {
		n.Cc = cc
	}

	return n, nil
}

// GetPostFromObjectURI attempt to derive a `Post` ID and its matching instance from an (ActivityPub) object URI.
func GetPostFromObjectURI(ctx context.Context, uris_table *uris.URIs, posts_db PostsDatabase, object_uri string) (*Post, error) {

	pat_post := uris_table.Post
	pat_post = strings.Replace(pat_post, "{resource}", "(?:@[^\\/]+)", 1)
	pat_post = strings.Replace(pat_post, "{id}", "(\\d+)", 1)

	re_post, err := regexp.Compile(pat_post)

	if err != nil {
		return nil, fmt.Errorf("Failed to compile post URI pattern, %w", err)
	}

	u, err := url.Parse(object_uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	object_path := u.Path

	if !re_post.MatchString(object_path) {
		slog.Debug("Invalid or unsupport post URI", "uri", object_uri)
		return nil, fmt.Errorf("Invalid or unsupport post URI")
	}

	m := re_post.FindStringSubmatch(object_path)

	str_id := m[1]
	post_id, err := strconv.ParseInt(str_id, 10, 64)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse ID derived from object URI, %w", err)
	}

	return posts_db.GetPostWithId(ctx, post_id)
}
