package utils

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type jwtManager struct{}

var (
	pubKey *rsa.PublicKey
	privKey *rsa.PrivateKey
	oncePubKey    sync.Once
	oncePrivKey    sync.Once
)

func (m *jwtManager) NewToken(username string) (string, error) {
	oncePrivKey.Do(func() {
		privKeyFile, err := os.ReadFile("vault/private-key.pem")
		if err != nil {
			panic(err)
		}

		privKey, err = jwt.ParseRSAPrivateKeyFromPEM(privKeyFile)
		if err != nil {
			panic(err)
		}
	})

	if privKey == nil {
		return "", errors.New("No token signature key")
	}

	tk := jwt.NewWithClaims(
		jwt.SigningMethodRS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Hour * 2).Unix(),
		},
	)

	tkStr, err := tk.SignedString(privKey)
	if err != nil {
		return "", err
	}

	return tkStr, nil
}

func (m *jwtManager) VerifyToken(tkStr string) error {
	oncePubKey.Do(func() {
		pubKeyFile, err := os.ReadFile("vault/public-key.pem")
		if err != nil {
			panic(err)
		}

		pubKey, err = jwt.ParseRSAPublicKeyFromPEM(pubKeyFile)
		if err != nil {
			panic(err)
		}
	})

	if pubKey == nil {
		return errors.New("No token verification key")
	}

	tk, err := jwt.Parse(tkStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %v", t.Header["alg"])
		}

		return pubKey, nil
	})

	if err != nil {
		return err
	}

	if !tk.Valid {
		return errors.New("Invalid token")
	}

	return nil
}

func NewJwtManager() *jwtManager {
	return &jwtManager{}
}
