package activitypub

import (
	"context"
	"fmt"
	"log/slog"
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

func (p *Post) Deliver(ctx context.Context, followers_db FollowersDatabase, q DeliveryQueue) error {

	followers_cb := func(ctx context.Context, follower_id string) error {

		slog.Info("Deliver", "post", p.Id, "follower_id", follower_id)
		err := q.DeliverPost(ctx, p, follower_id)

		if err != nil {
			return fmt.Errorf("Failed to deliver post to %s, %w", follower_id, err)
		}

		return nil
	}

	err := followers_db.GetFollowers(ctx, p.AccountId, followers_cb)

	if err != nil {
		return fmt.Errorf("Failed to get followers for post author, %w", err)
	}

	return nil
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

func (p *Post) AsCreateActivity(ctx context.Context, to []string) (*ap.Activity, error) {

	note, err := p.AsNote(ctx)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive note from post, %w", err)
	}

	create_activity, err := ap.NewCreateActivity(ctx, p.AccountId, to, note)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive create activity from note, %w", err)
	}

	return create_activity, nil
}
