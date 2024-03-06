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

type Post struct {
	Id        int64 `json:"id"`
	AccountId int64 `json:"account_id"`
	// This is a string mostly because []byte thingies get encoded incorrectly
	// in DynamoDB
	Body         string `json:"body"`
	InReplyTo    string `json:"in_reply_to"`
	Created      int64  `json:"created"`
	LastModified int64  `json:"lastmodified"`
}

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

func NoteFromPost(ctx context.Context, uris_table *uris.URIs, acct *Account, post *Post, post_tags []*PostTag) (*ap.Note, error) {

	attr := acct.ProfileURL(ctx, uris_table).String()
	post_url := acct.PostURL(ctx, uris_table, post)

	t := time.Unix(post.Created, 0)

	tags := make([]*ap.Tag, len(post_tags))

	for idx, pt := range post_tags {

		t := &ap.Tag{
			Name: pt.Name,
			Href: pt.Href,
			Type: pt.Type,
		}

		tags[idx] = t
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

	return n, nil
}

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
