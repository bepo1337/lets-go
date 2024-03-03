package mocks

import (
	"letsgo.bepo1337/internal/models"
	"time"
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
	if email == "alice@email" && password == "123" {
		return 1, nil
	} else {
		return -1, models.ErrInvalidCredentials
	}
}
func (u *UserModel) Get(id int) (*models.User, error) {
	if id == 1 {
		return &models.User{
			Id:      1,
			Name:    "Alice",
			Email:   "Random Mail",
			Created: time.Now(),
		}, nil
	} else {
		return &models.User{}, models.ErrNoRecord
	}
}

func (u *UserModel) CorrectPassword(id int, password string) (bool, error) {
	return true, nil
}

func (u *UserModel) UpdatePassword(id int, newPassword string) (bool, error) {
	return true, nil
}
