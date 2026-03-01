package repository

import (
	"context"
	"time"
)

type TokenRepository interface {
	SaveAccessToken(ctx context.Context, userID int64, token string, expiredAt time.Time) error

	IsAccessTokenValid(ctx context.Context, userID int64, token string) (bool, error)

	DeleteAccessToken(ctx context.Context, userID int64) error

	DeleteExpiredTokens(ctx context.Context) error

	BlacklistToken(ctx context.Context, userID int64, token string, expiredAt time.Time) error

	MoveToBlacklist(ctx context.Context, userID int64) error
}
