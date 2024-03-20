package storage

import "errors"

var (
	ErrNotFound      = errors.New("not found")
	ErrExists        = errors.New("already exists")
	ErrWrongPassword = errors.New("wrong password")
	ErrBlocked       = errors.New("user is blocked")
	ErrExpired       = errors.New("token expired")
)
