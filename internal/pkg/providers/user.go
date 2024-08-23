package providers

import (
	"errors"
	"texting-app/internal/pkg/models"
	"texting-app/internal/pkg/utils"
)

type userStore interface {
	GetUser(username string) (*models.User, error)
}

type UserProvider struct {
	store userStore
}

func (s *UserProvider) AuthUser(username string, password string) (*models.User, error) {
	user, err := s.store.GetUser(username)
	if err != nil {
		return nil, err
	}

	cryptMngr := utils.NewCryptoManager(
		"vault/public-key.pem",
		"vault/private-key.pem",
	)
	decrypted, err := cryptMngr.Decrypt(user.Password, username)
	if err != nil {
		return nil, err
	}

	if password != string(decrypted) {
		return nil, errors.New("unauthorized")
	}

	return user, nil
}

func NewUserProvider(s userStore) *UserProvider {
	return &UserProvider{
		store: s,
	}
}
