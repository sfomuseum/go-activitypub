package activitypub

import (
	"context"

	"github.com/sfomuseum/runtimevar"
)

type Actor struct {
	Id            string `json:"id"`
	PublicKeyURI  string `json:"public_key_uri"`
	PrivateKeyURI string `json:"private_key_uri"`
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
