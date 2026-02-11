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

func (r *ClientRepo) AddToCart(ctx context.Context, cart entity.Order) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("add to cart: %w", err)
	}
	defer tx.Rollback(ctx)

	stmt := `INSERT INTO cart (client_id) VALUES ($1) RETURNING cart_id`
	row := tx.QueryRow(ctx, stmt, cart.ClientID)
	if err := row.Scan(&cart.ID); err != nil {
		log.Println(err)
		return fmt.Errorf("add to cart: %w", err)
	}

	itemStmt := `INSERT INTO cart_items (cart_id, seller_id, product_id, quantity,price_at_purchase) VALUES ($1, $2, $3, $4,$5) RETURNING cart_item_id`
	for i := range cart.Items {
		item := &cart.Items[i]
		item.ID = cart.ID
		row = tx.QueryRow(ctx, itemStmt, item.ID, item.SellerID, item.ProductID, item.Quantity, item.PriceAtPurchase)
		if err := row.Scan(&item.ID); err != nil {
			log.Println(err)
			return fmt.Errorf("add to cart: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		log.Println(err)
		return fmt.Errorf("add to cart: %w", err)
	}

	return nil
}

func (r *ClientRepo) GetCartItems(ctx context.Context, clientID int64) ([]entity.OrderItem, error) {
	stmt := `SELECT ci.cart_item_id, ci.cart_id, ci.seller_id, ci.product_id, ci.quantity, ci.price_at_purchase
	FROM cart_items ci
	JOIN cart c ON ci.cart_id = c.cart_id
	WHERE c.client_id = $1`

	rows, err := r.conn.Query(ctx, stmt, clientID)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("get cart items: %w", err)
	}
	defer rows.Close()

	var items []entity.OrderItem
	for rows.Next() {
		var item entity.OrderItem
		if err := rows.Scan(&item.ID, &item.OrderID, &item.SellerID, &item.ProductID, &item.Quantity, &item.PriceAtPurchase); err != nil {
			log.Println(err)
			return nil, fmt.Errorf("get cart items: %w", err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, fmt.Errorf("get cart items: %w", err)
	}

	return items, nil
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

func (r *ClientRepo) AddAddress(ctx context.Context, client entity.Client) error {
	stmt := `UPDATE clients SET address=$1 WHERE client_id=$2`
	_, err := r.conn.Exec(ctx, stmt, client.Address, client.ID)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("add client address: %w", err)
	}
	return nil
}
