package models

import "errors"

var (
	ErrNoRecord           = errors.New("models: record with the id doesnt exist")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
)
