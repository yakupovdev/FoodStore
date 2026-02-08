package impl

import (
	"context"
	"fmt"
	"log"

	pg "github.com/jackc/pgx/v5"
	"github.com/yakupovdev/FoodStore/internal/domain/entity"
)

type OrderRepo struct {
	conn *pg.Conn
}

func NewOrderRepo(conn *pg.Conn) *OrderRepo {
	return &OrderRepo{conn: conn}
}

func (r *OrderRepo) FindByClientID(ctx context.Context, clientID int64) ([]entity.Order, error) {
	stmt := `SELECT order_id, client_id, status, created_at FROM orders WHERE client_id=$1`

	rows, err := r.conn.Query(ctx, stmt, clientID)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("get orders: %w", err)
	}
	defer rows.Close()

	var orders []entity.Order
	for rows.Next() {
		var order entity.Order
		if err := rows.Scan(&order.ID, &order.ClientID, &order.Status, &order.CreatedAt); err != nil {
			log.Println(err)
			return nil, fmt.Errorf("get orders: %w", err)
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, fmt.Errorf("get orders: %w", err)
	}

	return orders, nil
}

func (r *OrderRepo) FindItemsByOrderID(ctx context.Context, orderID int64) ([]entity.OrderItem, error) {
	stmt := `SELECT order_item_id, order_id, seller_id, product_id, quantity, price_at_purchase FROM orders_items WHERE order_id=$1`

	rows, err := r.conn.Query(ctx, stmt, orderID)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("get order items: %w", err)
	}
	defer rows.Close()

	var items []entity.OrderItem
	for rows.Next() {
		var item entity.OrderItem
		if err := rows.Scan(&item.ID, &item.OrderID, &item.SellerID, &item.ProductID, &item.Quantity, &item.PriceAtPurchase); err != nil {
			log.Println(err)
			return nil, fmt.Errorf("get order items: %w", err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, fmt.Errorf("get order items: %w", err)
	}

	return items, nil
}
