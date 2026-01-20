package repository

import (
	"context"
	"errors"
	"fmt"

	pg "github.com/jackc/pgx/v5"
)

func (p *Postgres) RegisterUser(email string, password string, userType string, balance int64) (int64, error) {
	ctx := context.Background()
	stmt := `INSERT INTO users (email, password_hash, type,created_at,last_enter,balance) VALUES ($1, $2, $3, NOW(), NOW(), $4) RETURNING userid`
	var userID int64
	if err := p.Conn.QueryRow(ctx, stmt, email, password, userType, balance).Scan(&userID); err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}
	return userID, nil
}

func (p *Postgres) LoginUser(email string, password string) (int64, error) {
	ctx := context.Background()
	stmt := `SELECT userid FROM users WHERE email=$1 AND password_hash=$2`

	var userID int64
	if err := p.Conn.QueryRow(ctx, stmt, email, password).Scan(&userID); err != nil {
		return 0, fmt.Errorf("failed to login user: %w", err)
	}
	stmt = `UPDATE users SET last_enter = NOW() WHERE email=$1 AND password_hash=$2;`
	_, err := p.Conn.Exec(ctx, stmt, email, password)
	if err != nil {
		return 0, fmt.Errorf("failed to update last_enter: %w", err)
	}

	return userID, nil
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
	stmt := `SELECT userid FROM users WHERE email=$1`

	var userID int64
	err := p.Conn.QueryRow(ctx, stmt, email).Scan(&userID)
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return 0, ErrUserNotFound
		} else {
			return 0, ErrQueryRow
		}
	}
	return userID, nil
}

func (p *Postgres) UpdateUserPassword(userID int64, newPassword string) error {
	ctx := context.Background()
	stmt := `UPDATE users SET password_hash=$1 WHERE userid=$2;`

	if _, err := p.Conn.Exec(ctx, stmt, newPassword, userID); err != nil {
		return ErrUpdatePassword
	}
	return nil
}
