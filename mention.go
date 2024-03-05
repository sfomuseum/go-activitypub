package activitypub

import (
	"context"
	"fmt"
	"time"

	"github.com/sfomuseum/go-activitypub/id"
)

func NewMention(ctx context.Context, post *Post, name string, href string) (*PostTag, error) {

	mention_id, err := id.NewId()

	if err != nil {
		return nil, fmt.Errorf("Failed to create new ID, %w", err)
	}

	now := time.Now()
	ts := now.Unix()

	t := &PostTag{
		Id:        mention_id,
		AccountId: post.AccountId,
		PostId:    post.Id,
		Name:      name,
		Href:      href,
		Type:      "Mention",
		Created:   ts,
	}

	return t, nil
}
