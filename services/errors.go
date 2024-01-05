package services

import "errors"

var (
	ErrEmailTaken    = errors.New("services: Email address is already in use")
	ErrUserNotFound  = errors.New("services: user could not be found")
	ErrWrongPassword = errors.New("authenticate: wrong password")
)
