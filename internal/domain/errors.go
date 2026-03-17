package domain

import "errors"

var (
	ErrNotUniqueLogin    = errors.New("login is not unique")
	ErrIncorrectPassword = errors.New("incorrect password")
	ErrPermissionDenied  = errors.New("permission denied")
	ErrNotFound          = errors.New("not found")
)
