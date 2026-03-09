package user

import "errors"

var (
	ErrUserExists        = errors.New("user already exists")
	ErrUserNotExists     = errors.New("user does not exists")
	ErrUserPasswordWrong = errors.New("user passwords not right or not correct")
)
