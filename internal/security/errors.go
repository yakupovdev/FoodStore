package security

import "errors"

var (
	ErrTokenGeneration     = errors.New("could not generate token")
	ErrInvalidRecoveryCode = errors.New("invalid recovery code")
)
