package models

type User struct {
    username string
	password string
}

func (u *User) GetToken() string {
	return ""
}
