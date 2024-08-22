package auth

import (
    "texting-app/internal/pkg/models"
    "texting-app/internal/pkg/repository"
)

func AuthenticateUser(username string, password string) (*models.User, error) {
    prov, err := repository.GetProvider()
    if err != nil {
        return nil, err
    }

    user, err := prov.GetUser(username)
    if err != nil{
        return nil, err
    }

    user

	return 
}
