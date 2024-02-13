package activitypub

import (
	"context"
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

func (a *Actor) WebfingerResource(uris_table *URIs) (*webfinger.Resource, error) {

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

func (a *Actor) ProfileResource(hostname string, uris_table *URIs) (*profile.Resource, error) {
	return nil, fmt.Errorf("Not implmented")
}

func (a *Actor) PublicKey(ctx context.Context) (string, error) {
	return a.loadRuntimeVar(ctx, a.PublicKeyURI)
}

func (a *Actor) PrivateKey(ctx context.Context) (string, error) {
	return a.loadRuntimeVar(ctx, a.PrivateKeyURI)
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
