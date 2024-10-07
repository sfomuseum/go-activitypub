package activitypub

import (
	"context"
	"fmt"
	"time"

	"github.com/sfomuseum/go-activitypub/id"
)

// Type Like is possibly (probably) a misnomer in the same way that type `Post` is (see notes in
// post.go and boost.go). Specifically this data and the correspinding `LikesDatabase` was created
// to record likes from external actors about posts created by accounts on this server. It is not
// currently suited to record or deliver likes of external posts made by accounts on this server.
type Like struct {
	Id        int64  `json:"id"`
	AccountId int64  `json:"account_id"`
	PostId    int64  `json:"post_id"`
	Actor     string `json:"actor"`
	Created   int64  `json:"created"`
}

func NewLike(ctx context.Context, post *Post, actor string) (*Like, error) {

	like_id, err := id.NewId()

	if err != nil {
		return nil, fmt.Errorf("Failed to create new ID, %w", err)
	}

	now := time.Now()
	ts := now.Unix()

	l := &Like{
		Id:        like_id,
		AccountId: post.AccountId,
		PostId:    post.Id,
		Actor:     actor,
		Created:   ts,
	}

	return l, nil
}
