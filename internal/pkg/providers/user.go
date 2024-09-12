package providers

import (
	"errors"
	"texting-app/internal/pkg/models"
	"texting-app/internal/pkg/utils"
)

type userStore interface {
	GetUser(usrn string) (*models.User, error)
	CreateUser(usrn, pw string) (*models.User, error)
	FindUser(forUrsn, usrn string, offset uint) ([]*models.User, error)
	GetUserFriends(usrn string, offset uint) ([]*models.User, error)
	SendFriendRequest(usrn, recUsrn string) error
	GetFriendRequests(usrn string, offset uint) ([]*models.FriendRequest, error)
}

type userProvider struct {
	store userStore
}

var cryptMngr = utils.NewCryptoManager(
	"vault/public-key.pem",
	"vault/private-key.pem",
)

func (p *userProvider) AuthUser(usrn, pw string) (*models.User, error) {
	usr, err := p.store.GetUser(usrn)
	if err != nil {
		return nil, err
	}

	decrypted, err := cryptMngr.Decrypt(usr.Password)
	if err != nil {
		return nil, err
	}

	if pw != string(decrypted) {
		return nil, errors.New("unauthorized")
	}

	return usr, nil
}

func (p *userProvider) GetUser(usrn string) (*models.User, error) {
	usr, err := p.store.GetUser(usrn)
	if err != nil {
		return nil, err
	}

	return usr, nil
}

func (p *userProvider) CreateUser(usrn, pw string) (*models.User, error) {
	encrypted, err := cryptMngr.Encrypt(pw)
	if err != nil {
		return nil, err
	}

	usr, err := p.store.CreateUser(usrn, encrypted)
	if err != nil {
		return nil, err
	}

	return usr, nil
}

func (p *userProvider) FindUser(forUsrn, usrn string, offset uint) ([]string, error) {
	usrs, err := p.store.FindUser(forUsrn, usrn, offset)
	if err != nil {
		return nil, err
	}

	usrns := make([]string, 0)
	for i := 0; i < len(usrs); i++ {
		usrns = append(usrns, usrs[i].Username)
	}

	return usrns, nil
}

func (p *userProvider) GetUserFriends(usrn string, offset uint) ([]*models.User, error) {
	return p.store.GetUserFriends(usrn, offset)
}

func (p *userProvider) SendFriendRequest(usrn, recUsrn string) error {
	return p.store.SendFriendRequest(usrn, recUsrn)
}

func (p *userProvider) GetFriendRequests(usrn string, offset uint) ([]*models.FriendRequest, error) {
	return p.store.GetFriendRequests(usrn, offset)
}

func NewUserProvider(s userStore) *userProvider {
	return &userProvider{
		store: s,
	}
}
