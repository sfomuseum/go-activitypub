package activitypub

import (
	"context"
	"fmt"
	"time"

	"github.com/sfomuseum/go-activitypub/id"
)

// Type Boosted represents an object/URI/thing that a sfomuseum/go-activitypub.Account
// has boosted. It remains TBD whether this should try to be merged with the `Boost`
// struct below which would really mean replace `Boost.PostId` with `Boost.Object`...
type Boosted struct {
	Id        int64  `json:"id"`
	AccountId int64  `json:"account_id"`
	Author    string `json:"author"`
	Object    string `json:"object"`
	Created   int64  `json:"created"`
}

// Type Boost is possibly (probably) a misnomer in the same way that type `Post` is (see notes in
// post.go). Specifically this data and the correspinding `BoostsDatabase` was created to record
// boosts from external actors about posts created by accounts on this server. It is not currently
// suited to record or deliver boosts of external posts made by accounts on this server.
type Boost struct {
	Id        int64  `json:"id"`
	AccountId int64  `json:"account_id"`
	PostId    int64  `json:"post_id"`
	Actor     string `json:"actor"`
	Created   int64  `json:"created"`
}

func NewBoost(ctx context.Context, post *Post, actor string) (*Boost, error) {

	boost_id, err := id.NewId()

	if err != nil {
		return nil, fmt.Errorf("Failed to create new ID, %w", err)
	}

	now := time.Now()
	ts := now.Unix()

	l := &Boost{
		Id:        boost_id,
		AccountId: post.AccountId,
		PostId:    post.Id,
		Actor:     actor,
		Created:   ts,
	}

	return l, nil
}
