package repository

import "errors"

var (
	ErrNoRecord       = errors.New("no matching record found")
	ErrDuplicateLogin = errors.New("duplicate login")
)
