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

	// Mastodon returns this for example

	other_pems := []string{
		"-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAtzewNBRhuKanadLIur8h\n02xvUYmHr58BzL2QDveRKdehk+EqYzOBIFQ+Vul9VoNdkTyUWeZ/C+o3Q8yTDRnU\nuF2XMnu6BiHrcZCR4+9enaiKO63K3PN5IwcvsQVapzjGLJaXYGv4i8pJA11SfD0K\n6oWNzjLza4Nw8V/G2nJU0rmUytsPNixzSyDjdl1k8JOPvdsMvQfHcdqBfLCFWE/I\nqDaBhskTk0HEcR2DtsKxCLVGpHTSt37BoucmmdKzcO7xziC981cM+ZAR5Q2JvvbR\nmRXzBcU5578DwVpM4Wfrm7SAZ6CqgZ38HJ39PBMmT+PyEbKprb5zh+/U9HLjERYt\nWQIDAQAB\n-----END PUBLIC KEY-----\n",
	}

	for idx, p := range other_pems {

		_, err = RSAPublicKeyFromPEM(p)

		if err != nil {
			t.Fatalf("Failed to parse other PEM at index %d, %v", idx, err)
		}
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
