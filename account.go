package activitypub

import (
	"context"
	"crypto/rsa"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/sfomuseum/go-activitypub/ap"
	"github.com/sfomuseum/go-activitypub/crypto"
	"github.com/sfomuseum/go-activitypub/uris"
	"github.com/sfomuseum/go-activitypub/webfinger"
	"github.com/sfomuseum/runtimevar"
)

type AccountType uint32

const (
	_ AccountType = iota
	PersonType
	ServiceType
)

func (t AccountType) String() string {

	switch t {
	case PersonType:
		return "Person"
	case ServiceType:
		return "Service"
	default:
		return ""
	}
}

func AccountTypeFromString(str_type string) (AccountType, error) {

	switch str_type {
	case "Person":
		return PersonType, nil
	case "Service":
		return ServiceType, nil
	default:
		return 0, fmt.Errorf("Invalid or unsupported account type")
	}
}

// https://www.w3.org/TR/activitypub/#actor-objects

type Account struct {
	Id                       int64       `json:"id"`
	AccountType              AccountType `json:"account_type"`
	Name                     string      `json:"name"`
	DisplayName              string      `json:"display_name"`
	Blurb                    string      `json:"blurb"`
	URL                      string      `json:"url"`
	PublicKeyURI             string      `json:"public_key_uri"`
	PrivateKeyURI            string      `json:"private_key_uri"`
	ManuallyApproveFollowers bool        `json:"manually_approve_followers"`
	Discoverable             bool        `json:"discoverable"`
	IconURI                  string      `json:"icon_uri"`
	Created                  int64       `json:"created"`
	LastModified             int64       `json:"lastmodified"`
}

func (a *Account) String() string {
	return a.Name
}

func (a *Account) Address(hostname string) string {
	return fmt.Sprintf("%s@%s", a.Name, hostname)
}

func (a *Account) AccountURL(ctx context.Context, uris_table *uris.URIs) *url.URL {

	account_path := uris.AssignResource(uris_table.Account, a.Name)
	return uris.NewURL(uris_table, account_path)
}

func (a *Account) OutboxURL(ctx context.Context, uris_table *uris.URIs) *url.URL {

	outbox_path := uris.AssignResource(uris_table.Outbox, a.Name)
	return uris.NewURL(uris_table, outbox_path)
}

func (a *Account) InboxURL(ctx context.Context, uris_table *uris.URIs) *url.URL {

	inbox_path := uris.AssignResource(uris_table.Inbox, a.Name)
	return uris.NewURL(uris_table, inbox_path)
}

func (a *Account) ProfileURL(ctx context.Context, uris_table *uris.URIs) *url.URL {

	account_path := uris.AssignResource(uris_table.Account, fmt.Sprintf("@%s", a.Name))
	return uris.NewURL(uris_table, account_path)
}

func (a *Account) PostURL(ctx context.Context, uris_table *uris.URIs, post *Post) *url.URL {

	account_path := uris.AssignResource(uris_table.Post, fmt.Sprintf("@%s", a.Name))
	post_path := uris.AssignId(account_path, strconv.FormatInt(post.Id, 10))
	return uris.NewURL(uris_table, post_path)
}

func (a *Account) WebfingerResource(ctx context.Context, uris_table *uris.URIs) (*webfinger.Resource, error) {

	account_url := a.AccountURL(ctx, uris_table)
	profile_url := a.ProfileURL(ctx, uris_table)

	subject := fmt.Sprintf("acct:%s@%s", a.Name, uris_table.Hostname)

	aliases := []string{
		account_url.String(),
		profile_url.String(),
	}

	profile_link := webfinger.Link{
		Rel:  "http://webfinger.net/rel/profile-page",
		Type: "text/html",
		HRef: account_url.String(),
	}

	activity_link := webfinger.Link{
		Rel:  "self",
		Type: "application/activity+json",
		HRef: account_url.String(),
	}

	links := []webfinger.Link{
		profile_link,
		activity_link,
	}

	r := &webfinger.Resource{
		Subject: subject,
		Aliases: aliases,
		Links:   links,
	}

	return r, nil
}

func (a *Account) FollowersResource(ctx context.Context, uris_table *uris.URIs, followers_database FollowersDatabase) (*ap.Followers, error) {

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

func (a *Account) FollowingResource(ctx context.Context, uris_table *uris.URIs, following_database FollowingDatabase) (*ap.Following, error) {

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

func (a *Account) ProfileResource(ctx context.Context, uris_table *uris.URIs) (*ap.Actor, error) {

	// https://www.w3.org/TR/activitypub/#actor-objects
	// https://www.w3.org/TR/activitypub/#obj-id
	// https://www.w3.org/TR/activitypub/#inbox
	// https://www.w3.org/TR/activitypub/#outbox
	// https://www.w3.org/TR/activitystreams-vocabulary/#dfn-orderedcollection

	account_url := a.AccountURL(ctx, uris_table)

	inbox_url := a.InboxURL(ctx, uris_table)
	outbox_url := a.OutboxURL(ctx, uris_table)

	icon_path := uris.AssignResource(uris_table.Icon, a.Name)
	icon_url := uris.NewURL(uris_table, icon_path)

	followers_path := uris.AssignResource(uris_table.Followers, a.Name)
	followers_url := uris.NewURL(uris_table, followers_path)

	following_path := uris.AssignResource(uris_table.Following, a.Name)
	following_url := uris.NewURL(uris_table, following_path)

	pem, err := runtimevar.StringVar(ctx, a.PublicKeyURI)

	if err != nil {
		return nil, fmt.Errorf("Failed to read public key URI, %w", err)
	}

	pub_key := ap.PublicKey{
		Id:    account_url.String() + "#main-key",
		Owner: account_url.String(),
		PEM:   pem,
	}

	icon := ap.Icon{
		Type:      "Image",
		MediaType: "image/png",
		URL:       icon_url.String(),
	}

	context := []interface{}{
		"https://www.w3.org/ns/activitystreams",
		"https://w3id.org/security/v1",
	}

	// read from prefs or something...
	discoverable := true
	// manually_approve := false

	now := time.Now()

	pr := &ap.Actor{
		Context:           context,
		Id:                account_url.String(),
		Type:              a.AccountType.String(),
		Name:              a.Name,
		PreferredUsername: a.Name,
		Summary:           a.Blurb,
		URL:               a.URL,
		Followers:         followers_url.String(),
		Following:         following_url.String(),
		// ManuallyApprovesFollowers: manually_approve,
		Discoverable: discoverable,
		Inbox:        inbox_url.String(),
		Outbox:       outbox_url.String(),
		PublicKey:    pub_key,
		Icon:         icon,
		Published:    now.Format(time.RFC3339),
	}

	return pr, nil
}

func (a *Account) PublicKey(ctx context.Context) (string, error) {
	return a.loadRuntimeVar(ctx, a.PublicKeyURI)
}

func (a *Account) PublicKeyRSA(ctx context.Context) (*rsa.PublicKey, error) {

	public_key_str, err := a.PublicKey(ctx)

	if err != nil {
		return nil, fmt.Errorf("Failed to get public key, %w", err)
	}

	return crypto.RSAPublicKeyFromPEM(public_key_str)
}

func (a *Account) PrivateKey(ctx context.Context) (string, error) {
	return a.loadRuntimeVar(ctx, a.PrivateKeyURI)
}

func (a *Account) PrivateKeyRSA(ctx context.Context) (*rsa.PrivateKey, error) {

	private_key_str, err := a.PrivateKey(ctx)

	if err != nil {
		return nil, fmt.Errorf("Failed to get private key, %w", err)
	}

	return crypto.RSAPrivateKeyFromPEM(private_key_str)
}

func (a *Account) loadRuntimeVar(ctx context.Context, uri string) (string, error) {
	return runtimevar.StringVar(ctx, uri)
}

func AddAccount(ctx context.Context, db AccountsDatabase, a *Account) (*Account, error) {

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

func UpdateAccount(ctx context.Context, db AccountsDatabase, a *Account) (*Account, error) {

	now := time.Now()
	ts := now.Unix()

	a.LastModified = ts

	err := db.UpdateAccount(ctx, a)

	if err != nil {
		return nil, fmt.Errorf("Failed to update account, %w", err)
	}

	return a, nil
}
