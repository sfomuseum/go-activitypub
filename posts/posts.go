package posts

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

	"github.com/sfomuseum/go-activitypub"
	"github.com/sfomuseum/go-activitypub/ap"
	"github.com/sfomuseum/go-activitypub/database"
	"github.com/sfomuseum/go-activitypub/uris"
)

type AddPostOptions struct {
	URIs             *uris.URIs
	PostsDatabase    database.PostsDatabase
	PostTagsDatabase database.PostTagsDatabase
}

// AddPost creates a new post record for 'body' and adds it to the post database. Then it parses 'body' looking
// for other ActivityPub addresses and records each as a "mention" in the post tags database. It returns the post
// and the list of post tags (mentions) for further processing as needed.
func AddPost(ctx context.Context, opts *AddPostOptions, acct *activitypub.Account, body string) (*activitypub.Post, []*activitypub.PostTag, error) {

	// Create the new post record

	p, err := activitypub.NewPost(ctx, acct, body)

	if err != nil {
		return nil, nil, fmt.Errorf("Failed to create new post, %w", err)
	}

	// Add the post to the database

	err = opts.PostsDatabase.AddPost(ctx, p)

	if err != nil {
		return nil, nil, fmt.Errorf("Failed to add post, %w", err)
	}

	// Determine other accounts mentioned in post

	addrs_mentioned, err := ap.ParseAddressesFromString(body)

	if err != nil {
		return nil, nil, fmt.Errorf("Failed to derive addresses mentioned in message body, %w", err)
	}

	// Create "mention" tags for any other accounts mentioned in post
	// Add each mention to the "post tags" database

	post_tags := make([]*activitypub.PostTag, 0)

	for _, name := range addrs_mentioned {

		actor, err := ap.RetrieveActor(ctx, name, opts.URIs.Insecure)

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

		t, err := activitypub.NewMention(ctx, p, mention_name, mention_href)

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

// ActivityFromPost should be kept in /ap but that causes Golang import cycle errors.
// ActivityFromPost should also be updated to accept media attachments.

// ActivityFromPost creates a new (ActivityPub) `Activity` instance derived from 'acct', 'post' and 'post_tags'.
func ActivityFromPost(ctx context.Context, uris_table *uris.URIs, acct *activitypub.Account, post *activitypub.Post, mentions []*activitypub.PostTag) (*ap.Activity, error) {

	from_u := acct.AccountURL(ctx, uris_table)
	from := from_u.String()

	logger := slog.Default()
	logger = logger.With("post id", post.Id)
	logger = logger.With("from", from)

	logger.Debug("Create note from post")

	note, err := NoteFromPost(ctx, uris_table, acct, post, mentions)

	if err != nil {
		logger.Error("Failed to create note from post", "error", err)
		return nil, fmt.Errorf("Failed to derive note from post, %w", err)
	}

	to := []string{
		ap.ACTIVITYSTREAMS_CONTEXT_PUBLIC,
	}

	logger.Debug("Create activity from note")

	// Something something something media attachments here...

	return ap.NewCreateActivity(ctx, uris_table, from, to, note)
}

// THIS SHOULD BE IN /ap BUT CAUSES IMPORT CYCLE ERRORS

// NoteFromPost creates a new (ActivityPub) `Note` instance derived from 'acct', 'post' and 'post_tags'.
func NoteFromPost(ctx context.Context, uris_table *uris.URIs, acct *activitypub.Account, post *activitypub.Post, post_tags []*activitypub.PostTag) (*ap.Note, error) {

	attr := acct.ProfileURL(ctx, uris_table).String()
	post_url := acct.PostURL(ctx, uris_table, post)

	t := time.Unix(post.Created, 0)

	tags := make([]*ap.Tag, len(post_tags))
	// cc := make([]string, 0)

	for idx, pt := range post_tags {

		t := &ap.Tag{
			Name: pt.Name,
			Href: pt.Href,
			Type: pt.Type,
		}

		tags[idx] = t

		if pt.Type == "Mention" {
			// cc = append(cc, pt.Name)
		}
	}

	n := &ap.Note{
		Type:         "Note",
		Id:           post_url.String(),
		AttributedTo: attr,
		To: []string{
			ap.ACTIVITYSTREAMS_CONTEXT_PUBLIC,
		},
		// Cc:        cc,
		Content:   post.Body,
		Published: t.Format(http.TimeFormat),
		InReplyTo: post.InReplyTo,
		Tags:      tags,
		URL:       post_url.String(),
	}

	return n, nil
}

// GetPostFromObjectURI attempt to derive a `Post` ID and its matching instance from an (ActivityPub) object URI.
func GetPostFromObjectURI(ctx context.Context, uris_table *uris.URIs, posts_db database.PostsDatabase, object_uri string) (*activitypub.Post, error) {

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
