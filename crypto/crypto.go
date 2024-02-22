package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"regexp"
)

var re_public_type = regexp.MustCompile(`^(?:RSA\s+)?PUBLIC KEY`)
var re_private_type = regexp.MustCompile(`^(?:RSA\s+)?PRIVATE KEY`)

func GenerateKeyPair(sz int) ([]byte, []byte, error) {

	key, err := rsa.GenerateKey(rand.Reader, sz)

	if err != nil {
		return nil, nil, fmt.Errorf("Failed to generate key, %w", err)
	}

	pub := key.Public()

	keyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		},
	)

	pubPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(pub.(*rsa.PublicKey)),
		},
	)

	return keyPEM, pubPEM, nil
}

func RSAPublicKeyFromPEM(str_pem string) (*rsa.PublicKey, error) {

	public_key_block, _ := pem.Decode([]byte(str_pem))

	if public_key_block == nil {
		return nil, fmt.Errorf("failed to decode PEM block containing public key, key block is nil")
	}

	if !re_public_type.MatchString(public_key_block.Type) {
		return nil, fmt.Errorf("Invalid or unknown key block type '%s'", public_key_block.Type)
	}

	public_key, err := x509.ParsePKCS1PublicKey(public_key_block.Bytes)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse PEM block containing public key, %w", err)
	}

	return public_key, nil
}

func RSAPrivateKeyFromPEM(str_pem string) (*rsa.PrivateKey, error) {

	private_key_block, _ := pem.Decode([]byte(str_pem))

	if private_key_block == nil {
		return nil, fmt.Errorf("failed to decode PEM block containing private key, key block is nil")
	}

	if !re_private_type.MatchString(private_key_block.Type) {
		return nil, fmt.Errorf("Invalid or unknown key block type '%s'", private_key_block.Type)
	}

	private_key, err := x509.ParsePKCS1PrivateKey(private_key_block.Bytes)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse private key, %w", err)
	}

	return private_key, nil
}
