package utils

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

type jwtManager struct{}

var (
	signKey *rsa.PrivateKey
	once    sync.Once
)

func (m *jwtManager) NewToken(username string) (string, error) {
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
		&jwt.SigningMethodRSA{},
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Hour * 2).Unix(),
		},
	)

	tokenString, err := token.SignedString(signKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (m *jwtManager) VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if signKey == nil {
			return nil, errors.New("No token verification key")
		}

		return signKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("Invalid token")
	}

	return nil
}

func NewJwtManager() *jwtManager {
	return &jwtManager{}
}
