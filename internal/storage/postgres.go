package storage

import (
	"context"
	"errors"
	"fmt"
	"log"

	pg "github.com/jackc/pgx/v5"
)

type Postgres struct {
	Conn *pg.Conn
}

func NewPostgres(conn *pg.Conn) *Postgres {
	return &Postgres{
		Conn: conn,
	}
}

func (p *Postgres) CreateUser(login string, password string, Type string, balance float64) (int, error) {
	log.Println("VATAHELL")
	ctx := context.Background()
	stmt := `INSERT INTO users (login, password, type,created_at,last_enter,balance) VALUES ($1, $2, $3, NOW(), NOW(), $4) RETURNING id`
	log.Println("SRABOTALO")
	var id int
	if err := p.Conn.QueryRow(ctx, stmt, login, password, Type, balance).Scan(&id); err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return id, nil
}

func (p *Postgres) LoginUser(login string, password string) (int, error) {
	ctx := context.Background()
	stmt := `SELECT id FROM users WHERE login=$1 AND password=$2`

	var id int
	if err := p.Conn.QueryRow(ctx, stmt, login, password).Scan(&id); err != nil {
		return 0, fmt.Errorf("failed to login user: %w", err)
	}
	stmt = `UPDATE users SET last_enter = NOW() WHERE login=$1 AND password=$2;`
	_, err := p.Conn.Exec(ctx, stmt, login, password)
	if err != nil {
		return 0, fmt.Errorf("failed to update last_enter: %w", err)
	}

	return id, nil
}

func (p *Postgres) GetUserByLogin(login string) (bool, error) {
	ctx := context.Background()
	stmt := `SELECT 1 FROM users WHERE login=$1`

	var one int
	err := p.Conn.QueryRow(ctx, stmt, login).Scan(&one)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, pg.ErrNoRows) {
		return false, nil
	}
	return false, fmt.Errorf("get user by login failed: %w", err)
}
