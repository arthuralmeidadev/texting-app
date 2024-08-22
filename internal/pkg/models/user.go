package models

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	Username string
	Password string
}

var key = []byte("19dyg-B2I79-QL1f9-ol1j9-n1D7YGH986Yzh7crT6Rg")
var (
	signKey *rsa.PrivateKey
	once    sync.Once
)

func (u *User) NewToken() (string, error) {
	once.Do(func() {
        signKeyFile, err := os.ReadFile("vault/signatureKey.pem")
		if err != nil {
			panic(err)
		}

		signKeyPEM, _ := pem.Decode(signKeyFile)
		signKey, err = x509.ParsePKCS1PrivateKey(signKeyPEM.Bytes)
		if err != nil {
			panic(err)
		}
	})

    if signKey == nil {
        return "", errors.New("No token signature key")
    }

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": u.Username,
			"exp":      time.Now().Add(time.Hour * 2).Unix(),
		},
	)

	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
