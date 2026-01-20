package usecase

import "errors"

var (
	ErrDuplicateEmail     = errors.New("user with this email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrTokenGeneration    = errors.New("could not generate token")
	ErrDatabaseConnection = errors.New("failed to connect to the database")
	ErrUpdatePassword     = errors.New("failed to update password")
	ErrUserNotFound       = errors.New("user not found")
	ErrCodeIsNotValid     = errors.New("code is not valid")
	ErrVerificationFailed = errors.New("verification failed")
)
