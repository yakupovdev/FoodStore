package impl

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	pg "github.com/jackc/pgx/v5"
)

type TokenRepo struct {
	conn *pg.Conn
}

func NewTokenRepo(conn *pg.Conn) *TokenRepo {
	return &TokenRepo{conn: conn}
}

func (r *TokenRepo) SaveAccessToken(ctx context.Context, userID int64, token string, expiredAt time.Time) error {
	stmt := `
INSERT INTO token_whitelist (userid, access_token_hash, expired_at)
VALUES ($1, $2, $3)
ON CONFLICT (userid,access_token_hash)
DO UPDATE SET
  access_token_hash = EXCLUDED.access_token_hash,
  expired_at = EXCLUDED.expired_at;
`
	if _, err := r.conn.Exec(ctx, stmt, userID, token, expiredAt); err != nil {
		log.Println(err)
		return fmt.Errorf("save access token: %w", err)
	}
	return nil
}

func (r *TokenRepo) IsAccessTokenValid(ctx context.Context, userID int64, token string) (bool, error) {
	stmt := `SELECT access_token_hash, expired_at FROM token_whitelist WHERE userid=$1`

	var storedToken string
	var expiredAt time.Time
	err := r.conn.QueryRow(ctx, stmt, userID).Scan(&storedToken, &expiredAt)
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("failed to verify access token: %w", err)
	}
	if storedToken != token || time.Now().After(expiredAt) {
		return false, nil
	}
	return true, nil
}

func (r *TokenRepo) DeleteAccessToken(ctx context.Context, userID int64) error {
	stmt := `DELETE FROM token_whitelist WHERE userid=$1;`

	if _, err := r.conn.Exec(ctx, stmt, userID); err != nil {
		return fmt.Errorf("failed to delete access token: %w", err)
	}
	return nil
}

func (r *TokenRepo) DeleteExpiredTokens(ctx context.Context) error {
	stmt := `DELETE FROM token_blacklist WHERE expired_at <= now();`
	if _, err := r.conn.Exec(ctx, stmt); err != nil {
		return fmt.Errorf("failed to delete expired access tokens: %w", err)
	}
	return nil
}

func (r *TokenRepo) BlacklistToken(ctx context.Context, userID int64, token string, expiredAt time.Time) error {
	stmt := `
INSERT INTO token_blacklist (userid, access_token_hash, expired_at)
VALUES ($1, $2, $3);
`
	if _, err := r.conn.Exec(ctx, stmt, userID, token, expiredAt); err != nil {
		return fmt.Errorf("failed to blacklist access token: %w", err)
	}
	return nil
}

func (r *TokenRepo) MoveToBlacklist(ctx context.Context, userID int64) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	insertStmt := `INSERT INTO token_blacklist (userid, access_token_hash, expired_at)
		SELECT userid, access_token_hash, expired_at
		FROM token_whitelist
		WHERE userid=$1`
	if _, err := tx.Exec(ctx, insertStmt, userID); err != nil {
		return fmt.Errorf("failed to insert token into blacklist: %w", err)
	}

	deleteStmt := `DELETE FROM token_whitelist WHERE userid=$1`
	if _, err := tx.Exec(ctx, deleteStmt, userID); err != nil {
		return fmt.Errorf("failed to delete token from whitelist: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
