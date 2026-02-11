package impl

import (
	"context"
	"errors"
	"fmt"

	pg "github.com/jackc/pgx/v5"
	"github.com/yakupovdev/FoodStore/internal/domain/entity"
)

type UserRepo struct {
	conn *pg.Conn
}

func NewUserRepo(conn *pg.Conn) *UserRepo {
	return &UserRepo{conn: conn}
}

func (r *UserRepo) Create(ctx context.Context, user *entity.User) (int64, error) {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var userID int64
	userStmt := `
        INSERT INTO users (email, password_hash, type, created_at, last_enter, balance) 
        VALUES ($1, $2, $3, NOW(), NOW(), $4) 
        RETURNING userid`

	err = tx.QueryRow(ctx, userStmt, user.Email, user.PasswordHash, user.UserType, user.Balance).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert into users: %w", err)
	}

	if user.UserType != "moderator" && user.UserType != "admin" {
		var secondaryStmt string
		switch user.UserType {
		case "client":
			secondaryStmt = `INSERT INTO clients (client_id, name) VALUES ($1, $2)`
		case "seller":
			secondaryStmt = `INSERT INTO sellers (seller_id, name) VALUES ($1, $2)`
		default:
			return 0, fmt.Errorf("unknown user type: %s", user.UserType)
		}

		_, err = tx.Exec(ctx, secondaryStmt, userID, user.Name)
		if err != nil {
			return 0, fmt.Errorf("failed to insert into %s: %w", user.UserType, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return userID, nil
}

func (r *UserRepo) FindByEmailAndType(ctx context.Context, email, userType string) (*entity.User, error) {
	stmt := `SELECT userid, email, password_hash, type, balance, created_at, last_enter
	         FROM users WHERE email=$1 AND type=$2`

	var u entity.User
	err := r.conn.QueryRow(ctx, stmt, email, userType).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.UserType, &u.Balance, &u.CreatedAt, &u.LastEnter,
	)
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("find user by email and type: %w", err)
	}

	return &u, nil
}

func (r *UserRepo) ExistsByEmailAndType(ctx context.Context, email, userType string) (bool, error) {
	stmt := `SELECT EXISTS(SELECT 1 FROM users WHERE email=$1 AND type=$2)`

	var exists bool
	err := r.conn.QueryRow(ctx, stmt, email, userType).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check user existence failed: %w", err)
	}

	return exists, nil
}

func (r *UserRepo) FindIDByEmailAndType(ctx context.Context, email, userType string) (int64, error) {
	stmt := `SELECT userid FROM users WHERE email=$1 AND type=$2`

	var userID int64
	err := r.conn.QueryRow(ctx, stmt, email, userType).Scan(&userID)
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return 0, fmt.Errorf("user not found: %w", err)
		}
		return 0, fmt.Errorf("find user id: %w", err)
	}
	return userID, nil
}

func (r *UserRepo) UpdatePassword(ctx context.Context, userID int64, passwordHash string) error {
	stmt := `UPDATE users SET password_hash=$1 WHERE userid=$2;`

	if _, err := r.conn.Exec(ctx, stmt, passwordHash, userID); err != nil {
		return fmt.Errorf("update password: %w", err)
	}
	return nil
}

func (r *UserRepo) UpdateLastLogin(ctx context.Context, userID int64) error {
	stmt := `UPDATE users SET last_enter = NOW() WHERE userid=$1`

	if _, err := r.conn.Exec(ctx, stmt, userID); err != nil {
		return fmt.Errorf("update last login: %w", err)
	}
	return nil
}

func (r *UserRepo) Delete(ctx context.Context, userID int64) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	deleteUserStmt := `DELETE FROM users WHERE userid=$1`
	if _, err := tx.Exec(ctx, deleteUserStmt, userID); err != nil {
		return fmt.Errorf("failed to delete from users: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (r *UserRepo) GetAllUsers(ctx context.Context) ([]entity.User, error) {
	stmt := `SELECT userid, email, type, balance, created_at, last_enter FROM users`

	rows, err := r.conn.Query(ctx, stmt)
	if err != nil {
		return nil, fmt.Errorf("query all users: %w", err)
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var u entity.User
		err := rows.Scan(&u.ID, &u.Email, &u.UserType, &u.Balance, &u.CreatedAt, &u.LastEnter)
		if err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate users: %w", err)
	}

	return users, nil
}
