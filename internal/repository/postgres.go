package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	pg "github.com/jackc/pgx/v5"
	"github.com/yakupovdev/FoodStore/internal/security"
)

type Postgres struct {
	Conn *pg.Conn
}

func NewPostgres(conn *pg.Conn) *Postgres {
	return &Postgres{
		Conn: conn,
	}
}

func (p *Postgres) CreateUser(email string, password string, Type string, balance float64) (int, error) {
	log.Println("VATAHELL")
	ctx := context.Background()
	stmt := `INSERT INTO users (email, password, type,created_at,last_enter,balance) VALUES ($1, $2, $3, NOW(), NOW(), $4) RETURNING id`
	log.Println("SRABOTALO")
	var id int
	if err := p.Conn.QueryRow(ctx, stmt, email, password, Type, balance).Scan(&id); err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return id, nil
}

func (p *Postgres) LoginUser(email string, password string) (int, error) {
	ctx := context.Background()
	stmt := `SELECT id FROM users WHERE email=$1 AND password=$2`

	var id int
	if err := p.Conn.QueryRow(ctx, stmt, email, password).Scan(&id); err != nil {
		return 0, fmt.Errorf("failed to login user: %w", err)
	}
	stmt = `UPDATE users SET last_enter = NOW() WHERE email=$1 AND password=$2;`
	_, err := p.Conn.Exec(ctx, stmt, email, password)
	if err != nil {
		return 0, fmt.Errorf("failed to update last_enter: %w", err)
	}

	return id, nil
}

func (p *Postgres) GetUserByEmail(email string) (bool, error) {
	ctx := context.Background()
	stmt := `SELECT 1 FROM users WHERE email=$1`

	var one int
	err := p.Conn.QueryRow(ctx, stmt, email).Scan(&one)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, pg.ErrNoRows) {
		return false, nil
	}
	return false, fmt.Errorf("get user by login failed: %w", err)
}

func (p *Postgres) GetUserIDByEmail(email string) (int64, error) {
	ctx := context.Background()
	stmt := `SELECT id FROM users WHERE email=$1`

	var id int64
	err := p.Conn.QueryRow(ctx, stmt, email).Scan(&id)
	fmt.Println(id)
	if err != nil {
		return 0, fmt.Errorf("get user ID by email failed: %w", err)
	}
	return id, nil
}

func (p *Postgres) SaveRecoveryCode(userID int64, email string, code string) error {
	ctx := context.Background()

	codeHash := security.HashPassword(code)
	expiredAt := time.Now().Add(10 * time.Minute)
	fmt.Println("AHUEEEEEET")

	stmt := `
INSERT INTO password_recovery_codes (user_id, email, code_hash, expired_at)
VALUES ($1, $2, $3, $4)
ON CONFLICT (user_id, email)
DO UPDATE SET
  code_hash = EXCLUDED.code_hash,
  expired_at = EXCLUDED.expired_at;
`
	if _, err := p.Conn.Exec(ctx, stmt, userID, email, codeHash, expiredAt); err != nil {
		return fmt.Errorf("failed to save recovery code: %w", err)
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

func (p *Postgres) UpdateUserPassword(userID int64, newPassword string) error {
	ctx := context.Background()
	stmt := `UPDATE users SET password=$1 WHERE id=$2;`

	if _, err := p.Conn.Exec(ctx, stmt, newPassword, userID); err != nil {
		return fmt.Errorf("failed to update user password: %w", err)
	}
	return nil
}
