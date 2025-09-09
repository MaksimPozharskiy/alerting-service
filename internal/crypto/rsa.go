package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

var (
	ErrFailedToReadKeyFile = errors.New("failed to read key file")
	ErrInvalidKey          = errors.New("invalid key: must be a PEM encoded PKCS1 or PKCS8 key")
	ErrNotPublicKey        = errors.New("provided key is not a public key")
	ErrNotPrivateKey       = errors.New("provided key is not a private key")
)

// LoadPublicKey loads a public key from a file.
func LoadPublicKey(filename string) (*rsa.PublicKey, error) {
	if filename == "" {
		return nil, nil
	}

	keyData, err := os.ReadFile(filename)
	if err != nil {
		return nil, ErrFailedToReadKeyFile
	}

	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, ErrInvalidKey
	}

	var pub interface{}
	pub, err = x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		pub, err = x509.ParsePKCS1PublicKey(block.Bytes)
		if err != nil {
			return nil, ErrInvalidKey
		}
	}

	publicKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, ErrNotPublicKey
	}

	return publicKey, nil
}

// LoadPrivateKey loads a private key from a file.
func LoadPrivateKey(filename string) (*rsa.PrivateKey, error) {
	if filename == "" {
		return nil, nil
	}

	keyData, err := os.ReadFile(filename)
	if err != nil {
		return nil, ErrFailedToReadKeyFile
	}

	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, ErrInvalidKey
	}

	var priv interface{}
	priv, err = x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		priv, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, ErrInvalidKey
		}
	}

	privateKey, ok := priv.(*rsa.PrivateKey)
	if !ok {
		return nil, ErrNotPrivateKey
	}

	return privateKey, nil
}

// EncryptData encrypts data using the provided public key.
func EncryptData(data []byte, publicKey *rsa.PublicKey) ([]byte, error) {
	if publicKey == nil {
		return data, nil
	}

	encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, data)
	if err != nil {
		return nil, err
	}
	return encrypted, nil
}

// DecryptData decrypts data using the provided private key.
func DecryptData(data []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	if privateKey == nil {
		return data, nil
	}
	decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, data)
	if err != nil {
		return nil, err
	}
	return decrypted, nil
}
