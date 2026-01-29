package repository

import "errors"

var (
	ErrNoRecord         = errors.New("no matching record found")
	ErrDuplicateLogin   = errors.New("duplicate login")
	ErrUserNotFound     = errors.New("user not found")
	ErrQueryRow         = errors.New("query row error")
	ErrSaveRecoveryCode = errors.New("save recovery code error")
	ErrUpdatePassword   = errors.New("update password error")
	ErrSaveAccessToken  = errors.New("save access token error")
)
