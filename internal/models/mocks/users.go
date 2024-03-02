package mocks

import (
	"letsgo.bepo1337/internal/models"
)

type UserModel struct{}

func (u *UserModel) Insert(name, email, password string) error {
	switch email {
	case "dupe":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (u *UserModel) Exists(id int) (bool, error) {
	switch id {
	case 1:
		return true, nil
	default:
		return false, nil
	}
}

func (u *UserModel) Authenticate(email, password string) (int, error) {
	if email == "alice" && password == "123" {
		return 1, nil
	} else {
		return -1, models.ErrInvalidCredentials
	}
}
