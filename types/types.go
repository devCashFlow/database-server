package types

import (
	"errors"
	"time"
)

type User struct {
	ID       int        `json:"id"`
	Email    string     `json:"email"`
	Name     string     `json:"name"`
	Surname  string     `json:"surname"`
	Username string     `json:"username"`
	Password *string    `json:"password,omitempty"`
	Created  *time.Time `json:"created,omitempty"`
	Role     *int       `json:"role,omitempty"`
}

func (u *User) ValidatePasswordHash(passHash string) (bool, error) {
	if u.Password == nil || *u.Password != passHash {
		return false, errors.New("Wrong Password")
	}
	return true, nil
}

type Email struct {
	Address string `json:"address"`
	Name    string `json:"name"`
}

type Contact struct {
	Name    string    `json:"name"`
	Message string    `json:"message"`
	Email   string    `json:"email"`
	Created time.Time `json:"created,omitempty"`
}
