package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"os"
)

type cryptoMng struct {
	pubKey  *rsa.PublicKey
	privKey *rsa.PrivateKey
}

func (c *cryptoMng) Encrypt(value string, factor string) ([]byte, error) {
    hash := sha256.New()
    hash.Write([]byte(factor))
    ciphertext, err := rsa.EncryptOAEP(
        hash,
        rand.Reader,
        c.pubKey,
        []byte(value),
        nil,
    )
    if err != nil {
        return nil, err
    }
	return ciphertext, nil
}

func (c *cryptoMng) Decrypt(ciphertext string, factor string) ([]byte, error) {
    hash := sha256.New()
    hash.Write([]byte(factor))
    plainText, err := rsa.DecryptOAEP(
        hash,
        rand.Reader,
        c.privKey,
        []byte(ciphertext),
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
func NewCryptoManager(pubKeyPath string, privKeyPath string) *cryptoMng {
	pubKeyFile, err := os.ReadFile(pubKeyPath)
	if err != nil {
		panic(err)
	}

	pubKeyPEM, _ := pem.Decode(pubKeyFile)
	pubKey, err := x509.ParsePKCS1PublicKey(pubKeyPEM.Bytes)
	if err != nil {
		panic(err)
	}

	privKeyFile, err := os.ReadFile(privKeyPath)
	if err != nil {
		panic(err)
	}

	privKeyPEM, _ := pem.Decode(privKeyFile)
	privKey, err := x509.ParsePKCS1PrivateKey(privKeyPEM.Bytes)
	if err != nil {
		panic(err)
	}

	return &cryptoMng{
		pubKey:  pubKey,
		privKey: privKey,
	}
}
