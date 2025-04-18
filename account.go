package activitypub

import (
	"context"
	"crypto/rsa"
	"fmt"
	"log/slog"
	"net/url"
	"strconv"
	"time"

	"github.com/sfomuseum/go-activitypub/ap"
	"github.com/sfomuseum/go-activitypub/crypto"
	"github.com/sfomuseum/go-activitypub/uris"
	"github.com/sfomuseum/go-activitypub/webfinger"
	"github.com/sfomuseum/runtimevar"
)

// AccountType denotes the ActivityPub account type
type AccountType uint32

const (
	_ AccountType = iota
	// PersonType is considered to be an actual human being
	PersonType
	// ServiceType is considered to be a "bot" or other automated account
	ServiceType
)

// String returns an English-language label for the AccountType
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

// AccountTypeFromString returns a known AccountType derived from 'str_type' (an English-language label)
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

// Account represents an individual ActivityPub account
type Account struct {
	// Id is a unique numeric identifier for the account
	Id int64 `json:"id"`
	// AccountType denotes the ActivityPub account type
	AccountType AccountType `json:"account_type"`
	// Name is the unique account name for the account
	Name string `json:"name"`
	// DisplayName is the long-form name for the account (which is not guaranteed to be unique across all accounts)
	DisplayName string `json:"display_name"`
	// Blurb is the descriptive text for the account
	Blurb string `json:"blurb"`
	// URL is the primary website associated with the account
	URL string `json:"url"`
	// PublicKeyURI is a valid `gocloud.dev/runtimevar` referencing the PEM-encoded public key for the account.
	PublicKeyURI string `json:"public_key_uri"`
	// PublicKeyURI is a valid `gocloud.dev/runtimevar` referencing the PEM-encoded private key for the account.
	PrivateKeyURI string `json:"private_key_uri"`
	// ManuallyApproveFollowers is a boolean flag signaling that follower requests need to be manually approved. Note: There are currently no tools or interfaces for handling those approvals.
	ManuallyApproveFollowers bool `json:"manually_approve_followers"`
	// Discoverable is a boolean flag signaling that the account is discoverable.
	Discoverable bool `json:"discoverable"`
	// IconURI is a valid `gocloud.dev/blob` URI (as in the bucket URI + filename) referencing the icon URI for the account.
	IconURI string `json:"icon_uri"`
	// Created is a Unix timestamp of when the account was created.
	Created int64 `json:"created"`
	// LastModified is a Unix timestamp of when the account was last modified.
	LastModified int64 `json:"lastmodified"`
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

func (a *Account) WebfingerURL(ctx context.Context, uris_table *uris.URIs) *url.URL {

	address := a.Address(uris_table.Hostname)
	acct := fmt.Sprintf("acct:%s", address)

	wf_q := &url.Values{}
	wf_q.Set("resource", acct)

	wf_u := &url.URL{}
	wf_u.Path = webfinger.Endpoint
	wf_u.RawQuery = wf_q.Encode()

	return wf_u
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
		Name:              a.DisplayName, // name is display name
		PreferredUsername: a.Name,        // preferred username is account (user)name
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

func (a *Account) SendActivity(ctx context.Context, uris_table *uris.URIs, inbox_uri string, activity *ap.Activity) error {

	logger := slog.Default()

	profile_url := a.AccountURL(ctx, uris_table)
	key_id := profile_url.String()

	logger = logger.With("account", a.String())
	logger = logger.With("inbox", inbox_uri)
	logger = logger.With("activity id", activity.Id)

	private_key, err := a.PrivateKeyRSA(ctx)

	if err != nil {
		logger.Error("Failed to derive private key", "error", err)
		return fmt.Errorf("Failed to derive private key for from account, %w", err)
	}

	logger.Debug("Post activity to inbox")

	return activity.PostToInbox(ctx, key_id, private_key, inbox_uri)
}
