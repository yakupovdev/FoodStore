package storage

import "errors"

var (
	ErrDatabaseConnection     = errors.New("failed to connect to the database")
	ErrDatabaseQuery          = errors.New("database query error")
	ErrDatabaseNotInitialized = errors.New("database not initialized")
)
