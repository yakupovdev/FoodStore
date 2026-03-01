package repository

import "context"

type RecoveryCodeRepository interface {
	Save(ctx context.Context, userID int64, email, userType, codeHash string, expiredAt interface{}) error

	Verify(ctx context.Context, email, userType, codeHash string) (bool, error)

	Delete(ctx context.Context, email, userType string) error

	DeleteExpired(ctx context.Context) error
}
