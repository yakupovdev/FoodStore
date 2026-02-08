package entity

import "time"

type TokenType string

const (
	AccessTokenType   TokenType = "access"
	RefreshTokenType  TokenType = "refresh"
	RecoveryTokenType TokenType = "recovery"
)

type TokenClaims struct {
	UserID   int64
	UserType string
}

type AccessToken struct {
	UserID    int64
	Token     string
	ExpiredAt time.Time
}

func (t *AccessToken) IsExpired() bool {
	return time.Now().After(t.ExpiredAt)
}
