package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"os"
)

type cryptoMng struct {
	pubKey  *rsa.PublicKey
	privKey *rsa.PrivateKey
}

func (c *cryptoMng) Encrypt(value string) (string, error) {
	ciphertext, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		c.pubKey,
		[]byte(value),
		nil,
	)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString([]byte(ciphertext)), nil
}

func (c *cryptoMng) Decrypt(ciphertext string) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, err
	}

	plainText, err := rsa.DecryptOAEP(
		sha256.New(),
		rand.Reader,
		c.privKey,
		decoded,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return plainText, nil
}

// Returns a pointer to a cryptoMng instance.
//
// If there is a problem reading the files or parsing them into RSA keys,
// calls panic() with the root error as argument
func NewCryptoManager(pubKeyPath, privKeyPath string) *cryptoMng {
	pubKeyFile, err := os.ReadFile(pubKeyPath)
	if err != nil {
		panic(err)
	}

	pubKeyPEM, _ := pem.Decode(pubKeyFile)
	pubKey, err := x509.ParsePKIXPublicKey(pubKeyPEM.Bytes)
	if err != nil {
		panic(err)
	}

	privKeyFile, err := os.ReadFile(privKeyPath)
	if err != nil {
		panic(err)
	}

	privKeyPEM, _ := pem.Decode(privKeyFile)
	privKey, err := x509.ParsePKCS8PrivateKey(privKeyPEM.Bytes)
	if err != nil {
		panic(err)
	}

	pubRsaKey, okPubKeyType := pubKey.(*rsa.PublicKey)
	privRsaKey, okPrivKeyType := privKey.(*rsa.PrivateKey)
	if !okPubKeyType {
		panic(errors.New("Invalid public key type. Expected RSA key"))
	} else if !okPrivKeyType {
		panic(errors.New("Invalid private key type. Expected RSA key"))
	}

	return &cryptoMng{
		pubKey:  pubRsaKey,
		privKey: privRsaKey,
	}
}
