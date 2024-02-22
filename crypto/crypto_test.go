package crypto

import (
	"testing"
)

func TestRSAPublicKeyFromPEM(t *testing.T) {

	_, pub, err := GenerateKeyPair(4096)

	if err != nil {
		t.Fatalf("Failed to generate key pair, %v", err)
	}

	_, err = RSAPublicKeyFromPEM(string(pub))

	if err != nil {
		t.Fatalf("Failed to derive public key, %v", err)
	}
}

func TestRSAPrivateKeyFromPEM(t *testing.T) {

	key, _, err := GenerateKeyPair(4096)

	if err != nil {
		t.Fatalf("Failed to generate key pair, %v", err)
	}

	_, err = RSAPrivateKeyFromPEM(string(key))

	if err != nil {
		t.Fatalf("Failed to derive private key, %v", err)
	}
}
