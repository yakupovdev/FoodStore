package impl

import (
	"context"
	"fmt"
	"log"

	pg "github.com/jackc/pgx/v5"
	"github.com/yakupovdev/FoodStore/internal/domain/entity"
)

type ClientRepo struct {
	conn *pg.Conn
}

func NewClientRepo(conn *pg.Conn) *ClientRepo {
	return &ClientRepo{conn: conn}
}

func (r *ClientRepo) FindByID(ctx context.Context, clientID int64) (*entity.Client, error) {
	stmt := `SELECT client_id, name, rating FROM clients WHERE client_id=$1`
	row := r.conn.QueryRow(ctx, stmt, clientID)

	var client entity.Client
	if err := row.Scan(&client.ID, &client.Name, &client.Rating); err != nil {
		log.Println(err)
		return nil, fmt.Errorf("get client profile: %w", err)
	}

	stmt = `SELECT email, balance FROM users WHERE userid=$1`
	row = r.conn.QueryRow(ctx, stmt, clientID)
	if err := row.Scan(&client.Email, &client.Balance); err != nil {
		log.Println(err)
		return nil, fmt.Errorf("get client profile: %w", err)
	}

	return &client, nil
}

func (r *ClientRepo) GetBalance(ctx context.Context, clientID int64) (int64, error) {
	stmt := `SELECT balance FROM users WHERE userid=$1`
	row := r.conn.QueryRow(ctx, stmt, clientID)

	var balance int64
	if err := row.Scan(&balance); err != nil {
		log.Println(err)
		return 0, fmt.Errorf("get client balance: %w", err)
	}

	return balance, nil
}

func (r *ClientRepo) UpdateBalance(ctx context.Context, clientID int64, newBalance int64) error {
	stmt := `UPDATE users SET balance = balance + $1 WHERE userid=$2`
	_, err := r.conn.Exec(ctx, stmt, newBalance, clientID)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("update client balance: %w", err)
	}
	return nil
}

func (r *ClientRepo) AddAddress(ctx context.Context, client entity.Client) error {
	stmt := `UPDATE clients SET address=$1 WHERE client_id=$2`
	_, err := r.conn.Exec(ctx, stmt, client.Address, client.ID)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("add client address: %w", err)
	}
	return nil
}
