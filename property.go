package activitypub

import (
	"context"
	"fmt"
	"time"

	"github.com/sfomuseum/go-activitypub/id"
)

// Property is a single key-value property associated with an account holder.
type Property struct {
	Id        int64  `json:"id"`
	AccountId int64  `json:"account_id"`
	Key       string `json:"key"`
	Value     string `json:"value"`
	Created   int64  `json:"created"`
}

func (pr *Property) String() string {
	return fmt.Sprintf("%s=%s", pr.Key, pr.Value)
}

func NewProperty(ctx context.Context, acct *Account, k string, v string) (*Property, error) {

	prop_id, err := id.NewId()

	if err != nil {
		return nil, err
	}

	now := time.Now()
	ts := now.Unix()

	pr := &Property{
		Id:        prop_id,
		AccountId: acct.Id,
		Key:       k,
		Value:     v,
		Created:   ts,
	}

	return pr, nil
}
