package activitypub

import (
	"context"
	"fmt"
	"time"

	"github.com/sfomuseum/go-activitypub/id"
)

// Note that we are not including the '//' part because that does not get
// included in serialized net/url.URL instances if the Host element is empty
const BOOST_URI_SCHEME string = "boost:"

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
