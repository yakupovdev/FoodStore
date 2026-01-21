package storage

import "errors"

var (
	ErrDatabaseConnection  = errors.New("failed to connect to the database")
	ErrUsersSchema         = errors.New("users schema error")
	ErrRecoveryCodesSchema = errors.New("recovery codes schema error")
)
