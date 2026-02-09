package impl

import (
	"context"

	pg "github.com/jackc/pgx/v5"
	"github.com/yakupovdev/FoodStore/internal/domain"
)

type TransactionRepository struct {
	conn *pg.Conn
}

func NewTransactionRepository(conn *pg.Conn) *TransactionRepository {
	return &TransactionRepository{conn: conn}
}

func (r *TransactionRepository) ExecuteOrderTransaction(ctx context.Context, clientID int64, totalAmount int64) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	isenoughBalanceStmt := `SELECT balance FROM users WHERE userid = $1`
	row := tx.QueryRow(ctx, isenoughBalanceStmt, clientID)

	var balance int64
	if err := row.Scan(&balance); err != nil {
		return err
	}

	if balance < totalAmount {
		return domain.ErrNotEnoughBalance
	}

	updateClientStmt := `UPDATE users SET balance = balance - $1 WHERE userid = $2`
	if _, err := tx.Exec(ctx, updateClientStmt, totalAmount, clientID); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *TransactionRepository) ExecuteSellerTransaction(ctx context.Context, sellerID int64, totalAmount int64) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	updateSellerStmt := `UPDATE users SET balance = balance + $1*0.95 WHERE userid = $2`
	if _, err := tx.Exec(ctx, updateSellerStmt, totalAmount, sellerID); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
