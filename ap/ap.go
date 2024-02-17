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

type Activity struct {
	Context string      `json:"@context"`
	Id      string      `json:"id"`
	Type    string      `json:"type"`
	Actor   string      `json:"actor"`
	To      []string    `json:"to,omitempty"`
	Object  interface{} `json:"object"`
}

type Note struct {
	Type         string      `json:"type"`
	Id           string      `json:"id"`
	AttributedTo string      `json:"attributedTo"`
	To           string      `json:"to"`
	Content      interface{} `json:"content"`
	URL          string      `json:"url"`
	Published    string      `json:"published"`
}

type PublicKey struct {
	Id    string `json:"id"`
	Owner string `json:"owner"`
	PEM   string `json:"publicKeyPem"`
}
