package activitypub

import (
	"context"
	"fmt"
	"time"

	"github.com/sfomuseum/go-activitypub/id"
)

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

func PropertiesMapForAccount(ctx context.Context, acct *Account, properties_db PropertiesDatabase) (map[string]*Property, error) {

	props_map := make(map[string]*Property)

	cb := func(ctx context.Context, pr *Property) error {

		_, exists := props_map[pr.Key]

		if exists {
			return fmt.Errorf("Duplicate key for %s", pr.Key)
		}

		props_map[pr.Key] = pr
		return nil
	}

	err := properties_db.GetPropertiesForAccount(ctx, acct.Id, cb)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive properties for account, %w", err)
	}

	return props_map, nil
}
