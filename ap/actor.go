package ap

import (
	"context"
	"crypto/rsa"
	"fmt"

	"github.com/sfomuseum/go-activitypub/crypto"
)

type Actor struct {
	Context           []string  `json:"@content"`
	Id                string    `json:"id"`
	Type              string    `json:"type"`
	PreferredUsername string    `json:"preferredUsername"`
	Inbox             string    `json:"inbox"`
	PublicKey         PublicKey `json:"publicKey"`

	Following                 string `json:"following,omitempty"`
	Followers                 string `json:"followers,omitempty"`
	Name                      string `json:"name,omitempty"`
	Summary                   string `json:"summary,omitempty"`
	URL                       string `json:"url,omitempty"`
	ManuallyApprovesFollowers bool   `json:"manuallyApprovesFollowers,omitempty"`
	Discoverable              bool   `json:"discoverable,omitempty"`
	Published                 string `json:"published,omitempty"`
	Icon                      Icon   `json:"icon,omitempty"`
}

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