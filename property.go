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

func PropertiesMapForAccount(ctx context.Context, properties_db PropertiesDatabase, acct *Account) (map[string]*Property, error) {

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

func DerivePropertiesUpdates(ctx context.Context, acct *Account, props_lookup map[string]*Property, updates map[string]string) ([]*Property, []*Property, error) {

	to_add := make([]*Property, 0)
	to_update := make([]*Property, 0)

	for k, v := range updates {

		pr, exists := props_lookup[k]

		if exists {

			if pr.Value != v {
				pr.Value = v
				to_update = append(to_update, pr)
			}

		} else {

			pr, err := NewProperty(ctx, acct, k, v)

			if err != nil {
				return nil, nil, fmt.Errorf("Failed to create %s property for %d, %w", k, acct.Id, err)
			}

			to_add = append(to_add, pr)
		}
	}

	return to_add, to_update, nil
}

func ApplyPropertiesUpdates(ctx context.Context, properties_db PropertiesDatabase, acct *Account, updates map[string]string) (int, int, error) {

	props_lookup, err := PropertiesMapForAccount(ctx, properties_db, acct)

	if err != nil {
		return 0, 0, fmt.Errorf("Failed to derive properties map for %d, %w", acct.Id, err)
	}

	to_add, to_update, err := DerivePropertiesUpdates(ctx, acct, props_lookup, updates)

	if err != nil {
		return 0, 0, fmt.Errorf("Failed to derive properties updates for %d, %w", acct.Id, err)
	}

	added := 0
	updated := 0

	for _, pr := range to_add {

		err := properties_db.AddProperty(ctx, pr)

		if err != nil {
			return added, updated, fmt.Errorf("Failed to add property (%s) for %d, %w", pr, acct.Id, err)
		}

		added += 1
	}

	for _, pr := range to_update {

		err := properties_db.UpdateProperty(ctx, pr)

		if err != nil {
			return added, updated, fmt.Errorf("Failed to update property (%s) for %d, %w", pr, acct.Id, err)
		}

		updated += 1
	}

	return added, updated, nil
}
