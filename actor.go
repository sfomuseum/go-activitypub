package activitypub

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/sfomuseum/go-activitypub/profile"
	"github.com/sfomuseum/go-activitypub/webfinger"
	"github.com/sfomuseum/runtimevar"
)

type Actor struct {
	Id            string `json:"id"`
	PublicKeyURI  string `json:"public_key_uri"`
	PrivateKeyURI string `json:"private_key_uri"`
	Created       int64  `json:"created"`
	LastModified  int64  `json:"lastmodified"`
}

func (a *Actor) String() string {
	return a.Id
}

func (a *Actor) ProfileURL(ctx context.Context, uris_table *URIs) (*url.URL, error) {

	id, hostname, err := ParseActorURI(a.Id)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse actor URI, %w", err)
	}

	profile_url := &url.URL{}
	profile_url.Scheme = "https"
	profile_url.Host = hostname
	profile_url.Path = filepath.Join(uris_table.Profile, id)

	return profile_url, nil
}

func (a *Actor) WebfingerResource(ctx context.Context, uris_table *URIs) (*webfinger.Resource, error) {

	id, hostname, err := ParseActorURI(a.Id)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse actor URI, %w", err)
	}

	subject := fmt.Sprintf("acct:%s", a.Id)

	profile_url := &url.URL{}
	profile_url.Scheme = "https"
	profile_url.Host = hostname
	profile_url.Path = filepath.Join(uris_table.Profile, id)

	activity_url := &url.URL{}
	activity_url.Scheme = "https"
	activity_url.Host = hostname
	activity_url.Path = filepath.Join(uris_table.Activity, id)

	profile_link := webfinger.Link{
		Rel:  "http://webfinger.net/rel/profile-page",
		Type: "text/html",
		HRef: profile_url.String(),
	}

	activity_link := webfinger.Link{
		Rel:  "self",
		Type: "application/activity+json",
		HRef: activity_url.String(),
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

func (a *Actor) ProfileResource(ctx context.Context, hostname string, uris_table *URIs) (*profile.Resource, error) {

	id, _, err := ParseActorURI(a.Id)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse actor URI, %w", err)
	}

	id_url := &url.URL{}
	id_url.Scheme = "https"
	id_url.Host = hostname
	id_url.Path = filepath.Join(uris_table.Id, id)

	inbox_url := &url.URL{}
	inbox_url.Scheme = "https"
	inbox_url.Host = hostname
	inbox_url.Path = filepath.Join(uris_table.Inbox, id)

	pem, err := runtimevar.StringVar(ctx, a.PublicKeyURI)

	if err != nil {
		return nil, fmt.Errorf("Failed to read public key URI, %w", err)
	}

	pub_key := profile.PublicKey{
		Id:    id_url.String() + "#main-key",
		Owner: id_url.String(),
		PEM:   pem,
	}

	context := []string{
		"https://www.w3.org/ns/activitystreams",
		"https://w3id.org/security/v1",
	}

	pr := &profile.Resource{
		Context:           context,
		Id:                id_url.String(),
		Type:              "Person",
		PreferredUsername: id,
		Inbox:             inbox_url.String(),
		PublicKey:         pub_key,
	}

	return pr, nil
}

func (a *Actor) PublicKey(ctx context.Context) (string, error) {
	return a.loadRuntimeVar(ctx, a.PublicKeyURI)
}

func (a *Actor) PrivateKey(ctx context.Context) (string, error) {
	return a.loadRuntimeVar(ctx, a.PrivateKeyURI)
}

func (a *Actor) PrivateKeyRSA(ctx context.Context) (*rsa.PrivateKey, error) {

	private_key_str, err := a.PrivateKey(ctx)

	if err != nil {
		return nil, fmt.Errorf("Failed to get private key, %w", err)
	}

	private_key_block, _ := pem.Decode([]byte(private_key_str))

	if err != nil {
		return nil, fmt.Errorf("Failed to decode private key, %w", err)
	}

	if private_key_block == nil || private_key_block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("Failed to decode PEM block containing private key")
	}

	private_key, err := x509.ParsePKCS1PrivateKey(private_key_block.Bytes)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse private key, %w", err)
	}

	return private_key, nil
}

func (a *Actor) loadRuntimeVar(ctx context.Context, uri string) (string, error) {
	return runtimevar.StringVar(ctx, uri)
}

func AddActor(ctx context.Context, db ActorDatabase, a *Actor) (*Actor, error) {

	now := time.Now()
	ts := now.Unix()

	a.Created = ts
	a.LastModified = ts

	err := db.AddActor(ctx, a)

	if err != nil {
		return nil, fmt.Errorf("Failed to add actor, %w", err)
	}

	return a, nil
}

func UpdateActor(ctx context.Context, db ActorDatabase, a *Actor) (*Actor, error) {

	now := time.Now()
	ts := now.Unix()

	a.LastModified = ts

	err := db.UpdateActor(ctx, a)

	if err != nil {
		return nil, fmt.Errorf("Failed to update actor, %w", err)
	}

	return a, nil
}

func ParseActorURI(uri string) (string, string, error) {

	parts := strings.Split(uri, "@")

	if len(parts) != 2 {
		return "", "", fmt.Errorf("Invalid address")
	}

	return parts[0], parts[1], nil
}
