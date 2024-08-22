package models

type User struct {
    Username string
	Password string
}

func (u *User) GetToken() string {
	return ""
}
