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
