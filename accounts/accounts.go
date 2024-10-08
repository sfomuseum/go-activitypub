package accounts

import (
	"context"
	"crypto/rsa"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/sfomuseum/go-activitypub"
	"github.com/sfomuseum/go-activitypub/ap"
	"github.com/sfomuseum/go-activitypub/database"	
	"github.com/sfomuseum/go-activitypub/crypto"
	"github.com/sfomuseum/go-activitypub/uris"
	"github.com/sfomuseum/go-activitypub/webfinger"
	"github.com/sfomuseum/runtimevar"
)

func AddAccount(ctx context.Context, db database.AccountsDatabase, a *activitypub.Account) (*activitypub.AAccount, error) {

	now := time.Now()
	ts := now.Unix()

	a.Created = ts
	a.LastModified = ts

	err := db.AddAccount(ctx, a)

	if err != nil {
		return nil, fmt.Errorf("Failed to add account, %w", err)
	}

	return a, nil
}

func UpdateAccount(ctx context.Context, db database.AccountsDatabase, a *activitypub.Account) (*activitypub.Account, error) {

	now := time.Now()
	ts := now.Unix()

	a.LastModified = ts

	err := db.UpdateAccount(ctx, a)

	if err != nil {
		return nil, fmt.Errorf("Failed to update account, %w", err)
	}

	return a, nil
}

func IsAccountNameTaken(ctx context.Context, db database.AccountsDatabase, name string) (bool, error) {

	_, err := db.GetAccountWithName(ctx, name)

	if err != nil {

		if err != ErrNotFound {
			return false, fmt.Errorf("Failed to determine is name is taken, %w", err)
		}

		return false, nil
	}

	return true, nil
}

func FollowersResource(ctx context.Context, uris_table *uris.URIs, a *activitypub.Account, followers_database database.FollowersDatabase) (*ap.Followers, error) {

	followers_path := uris.AssignResource(uris_table.Followers, a.Name)
	followers_url := uris.NewURL(uris_table, followers_path)

	count, err := CountFollowers(ctx, followers_database, a.Id)

	if err != nil {
		return nil, fmt.Errorf("Failed to count followers, %w", err)
	}

	f := &ap.Followers{
		Context:    "https://www.w3.org/ns/activitystreams",
		Id:         followers_url.String(),
		Type:       "OrderedCollection",
		TotalItems: count,
		First:      followers_url.String(),
	}

	return f, nil
}

func FollowingResource(ctx context.Context, uris_table *uris.URIs, a *activitypub.Account, following_database database.FollowingDatabase) (*ap.Following, error) {

	following_path := uris.AssignResource(uris_table.Following, a.Name)
	following_url := uris.NewURL(uris_table, following_path)

	count, err := CountFollowing(ctx, following_database, a.Id)

	if err != nil {
		return nil, fmt.Errorf("Failed to count following, %w", err)
	}

	f := &ap.Following{
		Context:    "https://www.w3.org/ns/activitystreams",
		Id:         following_url.String(),
		Type:       "OrderedCollection",
		TotalItems: count,
		First:      following_url.String(),
	}

	return f, nil
}
