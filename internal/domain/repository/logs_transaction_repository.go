package repository

import (
	"context"

	"github.com/yakupovdev/FoodStore/internal/domain/entity"
)

type LogsRepository interface {
	LogTransaction(ctx context.Context, log entity.LogTransaction) error

	GetLogsHistory(ctx context.Context) ([]entity.LogTransaction, error)
}
