package activitypub

import (
	"context"
	"fmt"
)

type Alias struct {
	Name      string `json:"name"` // Unique primary key
	AccountId int64  `json:"id"`
	Created   int64  `json:"created"`
}

func (a *Alias) String() string {
	return a.Name
}

func IsAliasNameTaken(ctx context.Context, aliases_db AliasesDatabase, name string) (bool, error) {

	_, err := aliases_db.GetAliasWithName(ctx, name)

	if err != nil {

		if err != ErrNotFound {
			return false, fmt.Errorf("Failed to determine is name is taken, %w", err)
		}

		return false, nil
	}

	return true, nil
}
