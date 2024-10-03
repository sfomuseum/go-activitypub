package ap

import (
	"context"
	"crypto/rsa"
	"fmt"

	"github.com/sfomuseum/go-activitypub/crypto"
)

// https://www.w3.org/TR/activitystreams-vocabulary/#actor-types

type Actor struct {
	// It has to be an interface because JSON-LD... thanks, JSON-LD...
	Context           []interface{} `json:"@context"`
	Id                string        `json:"id"`
	Type              string        `json:"type"`
	PreferredUsername string        `json:"preferredUsername"`
	Inbox             string        `json:"inbox"`
	Outbox            string        `json:"outbox"`
	PublicKey         PublicKey     `json:"publicKey"`

	Following string `json:"following,omitempty"`
	Followers string `json:"followers,omitempty"`
	Name      string `json:"name,omitempty"`
	Summary   string `json:"summary,omitempty"`
	URL       string `json:"url,omitempty"`
	// Don't omitempty because if you do then false values are omitted
	// ManuallyApprovesFollowers bool   `json:"manuallyApprovesFollowers"`
	Discoverable bool          `json:"discoverable,omitempty"`
	Published    string        `json:"published,omitempty"`
	Icon         Icon          `json:"icon,omitempty"`
	Attachments  []*Attachment `json:"attachment,omitempty"` // Is this just a Mastodon-ism?
}

// Returns the `rsa.PublicKey` instance for 'a'.
func (a *Actor) PublicKeyRSA(ctx context.Context) (*rsa.PublicKey, error) {

	public_key_str := a.PublicKey.PEM

	if public_key_str == "" {
		return nil, fmt.Errorf("Actor missing public key")
	}

	public_key, err := crypto.RSAPublicKeyFromPEM(public_key_str)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse PEM block containing public key, %w", err)

	}

	return public_key, nil
}
