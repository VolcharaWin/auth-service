package custom_errors

import "errors"

var (
	ErrLoginExists   = errors.New("this login already exists")
	ErrNotValidToken = errors.New("not a valid token")
)
