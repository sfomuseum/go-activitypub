package activitypub

import (
	"context"
	"fmt"
	"time"

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

func (a *Actor) WebfingerResource() (*webfinger.Resource, error) {

	subject := fmt.Sprintf("acct:%s", a.Id)

	profile_link := webfinger.Link{
		Rel:  "http://webfinger.net/rel/profile-page",
		Type: "text/html",
		HRef: fmt.Sprintf("/u/%s", a.Id),
	}

	activity_link := webfinger.Link{
		Rel:  "self",
		Type: "application/activity+json",
		HRef: fmt.Sprintf("/u/%s/activity", a.Id),
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
