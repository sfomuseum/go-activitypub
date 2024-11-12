package properties

import (
	"context"
	"fmt"

	"github.com/sfomuseum/go-activitypub"
	"github.com/sfomuseum/go-activitypub/database"
)

func PropertiesMapForAccount(ctx context.Context, properties_db database.PropertiesDatabase, acct *activitypub.Account) (map[string]*activitypub.Property, error) {

	props_map := make(map[string]*activitypub.Property)

	cb := func(ctx context.Context, pr *activitypub.Property) error {

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

func DerivePropertiesUpdates(ctx context.Context, acct *activitypub.Account, props_lookup map[string]*activitypub.Property, updates map[string]string) ([]*activitypub.Property, []*activitypub.Property, error) {

	to_add := make([]*activitypub.Property, 0)
	to_update := make([]*activitypub.Property, 0)

	for k, v := range updates {

		pr, exists := props_lookup[k]

		if exists {

			if pr.Value != v {
				pr.Value = v
				to_update = append(to_update, pr)
			}

		} else {

			pr, err := activitypub.NewProperty(ctx, acct, k, v)

			if err != nil {
				return nil, nil, fmt.Errorf("Failed to create %s property for %d, %w", k, acct.Id, err)
			}

			to_add = append(to_add, pr)
		}
	}

	return to_add, to_update, nil
}

func ApplyPropertiesUpdates(ctx context.Context, properties_db database.PropertiesDatabase, acct *activitypub.Account, updates map[string]string) (int, int, error) {

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
