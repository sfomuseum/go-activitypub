package activitypub

import (
	"context"
	"crypto/rsa"
	"fmt"
	"net/url"
	"time"

	"github.com/sfomuseum/go-activitypub/crypto"
	"github.com/sfomuseum/go-activitypub/profile"
	"github.com/sfomuseum/go-activitypub/webfinger"
	"github.com/sfomuseum/runtimevar"
)

type Account struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	PublicKeyURI  string `json:"public_key_uri"`
	PrivateKeyURI string `json:"private_key_uri"`
	Created       int64  `json:"created"`
	LastModified  int64  `json:"lastmodified"`
}

func (a *Account) String() string {
	return a.Name
}

func (a *Account) Address(hostname string) string {
	return fmt.Sprintf("%s@%s", a.Name, hostname)
}

func (a *Account) AccountURL(ctx context.Context, uris_table *URIs) *url.URL {

	account_path := AssignResource(uris_table.Account, a.Name)
	return NewURL(uris_table, account_path)
}

func (a *Account) WebfingerResource(ctx context.Context, uris_table *URIs) (*webfinger.Resource, error) {

	subject := fmt.Sprintf("acct:%s", a.Name)

	profile_url := a.AccountURL(ctx, uris_table)

	// activity_path := AssignResource(uris_table.Activity, a.Name)
	// activity_url := NewURL(uris_table, activity_path)

	profile_link := webfinger.Link{
		Rel:  "http://webfinger.net/rel/profile-page",
		Type: "text/html",
		HRef: profile_url.String(),
	}

	activity_link := webfinger.Link{
		Rel:  "self",
		Type: "application/activity+json",
		HRef: profile_url.String(),
	}

	links := []webfinger.Link{
		profile_link,
		activity_link,
	}

	r := &webfinger.Resource{
		Subject: subject,
		Links:   links,
	}

	return r, nil
}

func (a *Account) ProfileResource(ctx context.Context, uris_table *URIs) (*profile.Resource, error) {

	account_url := a.AccountURL(ctx, uris_table)

	inbox_path := AssignResource(uris_table.Inbox, a.Name)
	inbox_url := NewURL(uris_table, inbox_path)

	pem, err := runtimevar.StringVar(ctx, a.PublicKeyURI)

	if err != nil {
		return nil, fmt.Errorf("Failed to read public key URI, %w", err)
	}

	pub_key := profile.PublicKey{
		Id:    account_url.String() + "#main-key",
		Owner: account_url.String(),
		PEM:   pem,
	}

	context := []string{
		"https://www.w3.org/ns/activitystreams",
		"https://w3id.org/security/v1",
	}

	pr := &profile.Resource{
		Context:           context,
		Id:                account_url.String(),
		Type:              "Person",
		PreferredUsername: a.Name,
		Inbox:             inbox_url.String(),
		PublicKey:         pub_key,
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
