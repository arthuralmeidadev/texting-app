package models

import "texting-app/internal/pkg/utils"

type User struct {
	Username string
	Password string
}

func (u *User) NewToken() (string, error) {
	jwtMngr := utils.NewJwtManager()
	return jwtMngr.NewToken(u.Username)
}
