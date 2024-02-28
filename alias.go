package activitypub

import (
	"context"
	"fmt"

	"github.com/sfomuseum/go-activitypub/webfinger"
)

type Alias struct {
	Name      string `json:"name"` // Unique primary key
	AccountId int64  `json:"id"`
	Created   int64  `json:"created"`
}

func (a *Alias) String() string {
	return a.Name
}

func AppendAliasesToWebfingerResource(ctx context.Context, aliases_db AliasesDatabase, acct *Account, wf *webfinger.Resource) error {

	aliases := make([]string, 0)

	aliases_cb := func(ctx context.Context, alias *Alias) error {
		aliases = append(aliases, alias.Name)
		return nil
	}

	err := aliases_db.GetAliasesForAccount(ctx, acct.Id, aliases_cb)

	if err != nil {
		return fmt.Errorf("Failed to retrieve aliases for account, %w", err)
	}

	wf.Aliases = aliases
	return nil
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
