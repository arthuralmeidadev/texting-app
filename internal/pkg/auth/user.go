package auth

import (
	"errors"
	"texting-app/internal/pkg/models"
	"texting-app/internal/pkg/repository"
	"texting-app/internal/pkg/utils"
)

func AuthenticateUser(username string, password string) (*models.User, error) {
	prov, err := repository.GetProvider()
	if err != nil {
		return nil, err
	}

	user, err := prov.GetUser(username)
	if err != nil {
		return nil, err
	}

	cryptMngr := utils.NewCryptoManager()
	decrypted, err := cryptMngr.Decrypt(user.Password, username)
	if err != nil {
		return nil, err
	}

	if password != *decrypted {
        return nil, errors.New("unauthorized") 
	}

	return nil, nil
}
