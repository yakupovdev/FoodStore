package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	pg "github.com/jackc/pgx/v5"
)

type RecoveryCodeRepo struct {
	conn *pg.Conn
}

func NewRecoveryCodeRepo(conn *pg.Conn) *RecoveryCodeRepo {
	return &RecoveryCodeRepo{conn: conn}
}

func (r *RecoveryCodeRepo) Save(ctx context.Context, userID int64, email, userType, codeHash string, _ interface{}) error {
	expiredAt := time.Now().Add(10 * time.Minute)
	stmt := `
INSERT INTO password_recovery_codes (userid, email, type, code_hash, expired_at)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (userid)
DO UPDATE SET
  code_hash = EXCLUDED.code_hash,
  expired_at = EXCLUDED.expired_at;
`
	if _, err := r.conn.Exec(ctx, stmt, userID, email, userType, codeHash, expiredAt); err != nil {
		return fmt.Errorf("save recovery code: %w", err)
	}
	return nil
}

func (r *RecoveryCodeRepo) Verify(ctx context.Context, email, userType, codeHash string) (bool, error) {
	stmt := `SELECT code_hash, expired_at FROM password_recovery_codes WHERE email=$1 AND type=$2`

	var storedHash string
	var expiredAt time.Time
	err := r.conn.QueryRow(ctx, stmt, email, userType).Scan(&storedHash, &expiredAt)
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("failed to verify recovery code: %w", err)
	}

	if time.Now().After(expiredAt) {
		return false, nil
	}

	if storedHash != codeHash {
		log.Println("code hash mismatch")
		return false, nil
	}

	return true, nil
}

func (r *RecoveryCodeRepo) Delete(ctx context.Context, email, userType string) error {
	stmt := `DELETE FROM password_recovery_codes WHERE email=$1 AND type=$2`

	if _, err := r.conn.Exec(ctx, stmt, email, userType); err != nil {
		return fmt.Errorf("failed to delete recovery code: %w", err)
	}
	return nil
}

func (r *RecoveryCodeRepo) DeleteExpired(ctx context.Context) error {
	stmt := `DELETE FROM password_recovery_codes WHERE expired_at <= now();`

	if _, err := r.conn.Exec(ctx, stmt); err != nil {
		return fmt.Errorf("failed to delete expired recovery codes: %w", err)
	}
	return nil
}
