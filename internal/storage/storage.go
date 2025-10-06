package storage

import "errors"

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrUserExists    = errors.New("user exists")
	ErrNothingUpdate = errors.New("nothing to update")
)
