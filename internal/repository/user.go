package repository

import (
	"context"
	"errors"
	"fmt"
	"log"

	pg "github.com/jackc/pgx/v5"
)

func (p *Postgres) RegisterUser(email string, password string, userType string, balance int64, name string) (int64, error) {
	ctx := context.Background()

	tx, err := p.Conn.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var userID int64
	userStmt := `
        INSERT INTO users (email, password_hash, type, created_at, last_enter, balance) 
        VALUES ($1, $2, $3, NOW(), NOW(), $4) 
        RETURNING userid`

	err = tx.QueryRow(ctx, userStmt, email, password, userType, balance).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert into users: %w", err)
	}

	var secondaryStmt string
	switch userType {
	case "client":
		secondaryStmt = `INSERT INTO clients (client_id, name) VALUES ($1, $2)`
	case "seller":
		secondaryStmt = `INSERT INTO sellers (seller_id, name) VALUES ($1, $2)`
	default:
		return 0, fmt.Errorf("unknown user type: %s", userType)
	}

	_, err = tx.Exec(ctx, secondaryStmt, userID, name)
	if err != nil {
		return 0, fmt.Errorf("failed to insert into %s: %w", userType, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return userID, nil
}

func (p *Postgres) LoginUser(email, password, userType string) (int64, error) {
	ctx := context.Background()
	stmt := `SELECT userid FROM users WHERE email=$1 AND password_hash=$2 AND type=$3`

	var userID int64
	if err := p.Conn.QueryRow(ctx, stmt, email, password, userType).Scan(&userID); err != nil {
		return 0, fmt.Errorf("failed to login user: %w", err)
	}
	stmt = `UPDATE users SET last_enter = NOW() WHERE email=$1 AND password_hash=$2;`
	_, err := p.Conn.Exec(ctx, stmt, email, password)
	if err != nil {
		return 0, fmt.Errorf("failed to update last_enter: %w", err)
	}

	return userID, nil
}

func (p *Postgres) UserExists(email, userType string) (bool, error) {
	ctx := context.Background()
	stmt := `SELECT EXISTS(SELECT 1 FROM users WHERE email=$1 AND type=$2)`

	var exists bool
	err := p.Conn.QueryRow(ctx, stmt, email, userType).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("check user existence failed: %w", err)
	}

	return exists, nil
}

func (p *Postgres) GetUserIDByEmailAndType(email, userType string) (int64, error) {
	ctx := context.Background()
	stmt := `SELECT userid FROM users WHERE email=$1 AND type=$2`

	var userID int64
	err := p.Conn.QueryRow(ctx, stmt, email, userType).Scan(&userID)
	if err != nil {
		log.Println(err)
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
