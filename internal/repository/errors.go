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
	ErrGetOrders        = errors.New("get orders error")
	ErrGetOrderItems    = errors.New("get order items error")
	ErrGetProfile       = errors.New("get profile error")
	ErrGetCategories    = errors.New("get categories error")
	ErrGetProducts      = errors.New("get products error")
)
