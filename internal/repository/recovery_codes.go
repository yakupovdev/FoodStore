package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	pg "github.com/jackc/pgx/v5"
	"github.com/yakupovdev/FoodStore/internal/security"
)

func (p *Postgres) SaveRecoveryCode(userID int64, email string, code string) error {
	ctx := context.Background()

	codeHash := security.HashPassword(code)
	expiredAt := time.Now().Add(10 * time.Minute)

	stmt := `
INSERT INTO password_recovery_codes (userid, email, code_hash, expired_at)
VALUES ($1, $2, $3, $4)
ON CONFLICT (userID, email)
DO UPDATE SET
  code_hash = EXCLUDED.code_hash,
  expired_at = EXCLUDED.expired_at;
`
	if _, err := p.Conn.Exec(ctx, stmt, userID, email, codeHash, expiredAt); err != nil {
		return ErrSaveRecoveryCode
	}
	return nil
}

func (p *Postgres) DeleteExpiredRecoveryCodes() error {
	ctx := context.Background()
	stmt := `DELETE FROM password_recovery_codes WHERE expired_at <= now();`

	if _, err := p.Conn.Exec(ctx, stmt); err != nil {
		return fmt.Errorf("failed to delete expired recovery codes: %w", err)
	}
	return nil
}

func (p *Postgres) VerifyRecoveryCode(email string, code string) (bool, error) {
	ctx := context.Background()
	stmt := `SELECT code_hash, expired_at FROM password_recovery_codes WHERE email=$1`

	var codeHash string
	var expiredAt time.Time
	err := p.Conn.QueryRow(ctx, stmt, email).Scan(&codeHash, &expiredAt)
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("failed to verify recovery code: %w", err)
	}
	code = security.HashPassword(code)

	if time.Now().After(expiredAt) {
		return false, nil
	}

	if codeHash != code {
		return false, nil
	}

	return true, nil
}

func (p *Postgres) DeleteRecoveryCode(email string) error {
	ctx := context.Background()
	stmt := `DELETE FROM password_recovery_codes WHERE email=$1;`

	if _, err := p.Conn.Exec(ctx, stmt, email); err != nil {
		return fmt.Errorf("failed to delete recovery code: %w", err)
	}
	return nil
}
