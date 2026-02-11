package impl

import (
	"context"
	"fmt"
	"log"

	pg "github.com/jackc/pgx/v5"
)

type AdminRepository struct {
	conn *pg.Conn
}

func NewAdminRepository(conn *pg.Conn) *AdminRepository {
	return &AdminRepository{conn: conn}
}

func (r *AdminRepository) UpdateBalance(ctx context.Context, userID int64, newBalance int64) error {
	stmt := `UPDATE users SET balance = balance + $1 WHERE userid=$2`
	_, err := r.conn.Exec(ctx, stmt, newBalance, userID)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("update client balance: %w", err)
	}
	return nil
}

func (r *AdminRepository) GetExpiringSubscriptions(ctx context.Context) ([]int64, error) {
	stmt := `SELECT seller_id FROM subscriptions WHERE created_at < NOW() - INTERVAL '30 days'`
	rows, err := r.conn.Query(ctx, stmt)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("get expiring subscriptions: %w", err)
	}
	defer rows.Close()

	var sellerIDs []int64
	for rows.Next() {
		var sellerID int64
		if err := rows.Scan(&sellerID); err != nil {
			log.Println(err)
			return nil, fmt.Errorf("get expiring subscriptions: %w", err)
		}
		sellerIDs = append(sellerIDs, sellerID)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, fmt.Errorf("get expiring subscriptions: %w", err)
	}

	return sellerIDs, nil
}

func (r *AdminRepository) CancelSubscription(ctx context.Context, sellerID int64) error {
	stmt := `DELETE FROM subscriptions WHERE seller_id=$1`
	_, err := r.conn.Exec(ctx, stmt, sellerID)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("cancel subscription: %w", err)
	}
	return nil
}

func (r *AdminRepository) SetPriorityToFalse(ctx context.Context, sellerID int64) error {
	stmt := `UPDATE sellers SET priority = 0 WHERE seller_id=$1`
	_, err := r.conn.Exec(ctx, stmt, sellerID)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("set priority to false: %w", err)
	}
	return nil
}
