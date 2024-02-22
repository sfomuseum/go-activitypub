package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

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

	var public_key any
	var err error

	switch public_key_block.Type {
	case "RSA PUBLIC KEY":
		public_key, err = x509.ParsePKCS1PublicKey(public_key_block.Bytes)
	case "PUBLIC KEY":
		public_key, err = x509.ParsePKIXPublicKey(public_key_block.Bytes)
	default:
		err = fmt.Errorf("Invalid or unknown key block type, %s", public_key_block.Type)
	}

	if err != nil {
		return nil, fmt.Errorf("Failed to parse PEM block containing public key, %w", err)
	}

	return public_key.(*rsa.PublicKey), nil
}

func RSAPrivateKeyFromPEM(str_pem string) (*rsa.PrivateKey, error) {

	private_key_block, _ := pem.Decode([]byte(str_pem))

	if private_key_block == nil {
		return nil, fmt.Errorf("failed to decode PEM block containing private key, key block is nil")
	}

	var private_key any
	var err error

	switch private_key_block.Type {
	case "RSA PRIVATE KEY":
		private_key, err = x509.ParsePKCS1PrivateKey(private_key_block.Bytes)
	case "PRIVATE KEY":
		private_key, err = x509.ParsePKCS8PrivateKey(private_key_block.Bytes)
	default:
		err = fmt.Errorf("Invalid or unknown key block type, %s", private_key_block.Type)
	}

	if err != nil {
		return nil, fmt.Errorf("Failed to parse private key, %w", err)
	}

	return private_key.(*rsa.PrivateKey), nil
}
