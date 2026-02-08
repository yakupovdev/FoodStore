package domain

import "errors"

var (
	ErrEmptyEmail      = errors.New("email must not be empty")
	ErrEmptyPassword   = errors.New("password must not be empty")
	ErrInvalidUserType = errors.New("user type must be 'client' or 'seller'")
	ErrEmptyName       = errors.New("name must not be empty")
)

var (
	ErrInvalidPrice      = errors.New("price must be positive")
	ErrInvalidQuantity   = errors.New("quantity must be positive")
	ErrEmptyProductName  = errors.New("product name must not be empty")
	ErrEmptyCategoryName = errors.New("category name must not be empty")
	ErrEmptySubCatName   = errors.New("sub-category name must not be empty")
	ErrCategoryID        = errors.New("invalid category id")
	ErrSubCategoryID     = errors.New("invalid sub-category id")
	ErrProductID         = errors.New("invalid product id")
	ErrSellerID          = errors.New("invalid seller id")
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user with this email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

var (
	ErrTokenGeneration = errors.New("failed to generate token")
	ErrTokenStorage    = errors.New("failed to store token")
	ErrTokenExpired    = errors.New("token has expired")
	ErrTokenInvalid    = errors.New("token is invalid or revoked")
	ErrTokenCleanup    = errors.New("failed to clean up expired tokens")
)

var (
	ErrCodeInvalid        = errors.New("recovery code is invalid")
	ErrCodeExpired        = errors.New("recovery code has expired")
	ErrVerificationFailed = errors.New("verification failed")
	ErrUpdatePassword     = errors.New("failed to update password")
)

var (
	ErrInternal           = errors.New("internal server error")
	ErrDatabaseConnection = errors.New("database connection failed")
	ErrSMTPFailed         = errors.New("failed to send email")
	ErrPasswordHash       = errors.New("failed to hash password")
)
