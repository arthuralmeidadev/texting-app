package providers

import (
	"errors"
	"texting-app/internal/pkg/models"
	"texting-app/internal/pkg/utils"
)

type userStore interface {
	GetUser(usrn string) (*models.User, error)
	CreateUser(usrn, pw string) (*models.User, error)
}

type UserProvider struct {
	store userStore
}

var cryptMngr = utils.NewCryptoManager(
	"vault/public-key.pem",
	"vault/private-key.pem",
)

func (p *UserProvider) AuthUser(usrn, pw string) (*models.User, error) {
	usr, err := p.store.GetUser(usrn)
	if err != nil {
		return nil, err
	}

	decrypted, err := cryptMngr.Decrypt(usr.Password, usrn)
	if err != nil {
		return nil, err
	}

	if pw != string(decrypted) {
		return nil, errors.New("unauthorized")
	}

	return usr, nil
}

func (p *UserProvider) GetUser(usrn string) (*models.User, error) {
	usr, err := p.store.GetUser(usrn)
	if err != nil {
		return nil, err
	}

	return usr, nil
}

func (p *UserProvider) CreateUser(usrn, pw string) (*models.User, error) {
	encrypted, err := cryptMngr.Encrypt(pw, usrn)
	if err != nil {
		return nil, err
	}

	usr, err := p.store.CreateUser(usrn, string(encrypted))
	if err != nil {
		return nil, err
	}

	return usr, nil
}

func NewUserProvider(s userStore) *UserProvider {
	return &UserProvider{
		store: s,
	}
}
