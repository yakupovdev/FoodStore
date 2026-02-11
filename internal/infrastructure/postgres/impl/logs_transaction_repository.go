package impl

import (
	"context"

	pg "github.com/jackc/pgx/v5"
	"github.com/yakupovdev/FoodStore/internal/domain/entity"
)

type LogsRepository struct {
	conn *pg.Conn
}

func NewLogsRepository(conn *pg.Conn) *LogsRepository {
	return &LogsRepository{conn: conn}
}

func (r *LogsRepository) LogTransaction(ctx context.Context, log entity.LogTransaction) error {
	stmt := `INSERT INTO logs_transactions (client_id, seller_id, total_amount, commission_amount, created_at) VALUES ($1, $2, $3, $4, NOW())`
	_, err := r.conn.Exec(ctx, stmt, log.ClientID, log.SellerID, log.TotalAmount, log.CommissionAmount)
	return err
}

func (r *LogsRepository) GetLogsHistory(ctx context.Context) ([]entity.LogTransaction, error) {
	stmt := `SELECT log_id, client_id, seller_id, total_amount, commission_amount, created_at FROM logs_transactions ORDER BY created_at DESC`
	rows, err := r.conn.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []entity.LogTransaction
	for rows.Next() {
		var log entity.LogTransaction
		if err := rows.Scan(&log.ID, &log.ClientID, &log.SellerID, &log.TotalAmount, &log.CommissionAmount, &log.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}
