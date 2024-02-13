package activitypub

import (
	"context"
	"fmt"
	"time"

	"github.com/sfomuseum/runtimevar"
)

type Actor struct {
	Id            string `json:"id"`
	PublicKeyURI  string `json:"public_key_uri"`
	PrivateKeyURI string `json:"private_key_uri"`
	Created       int64  `json:"created"`
	LastModified  int64  `json:"lastmodified"`
}

func (a *Actor) Webfinger() (*Webfinger, error) {
	wf := &Webfinger{}
	return wf, nil
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
