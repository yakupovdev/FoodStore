package repository

import (
	"context"

	"github.com/yakupovdev/FoodStore/internal/domain/entity"
)

type OrderRepository interface {
	FindByClientID(ctx context.Context, clientID int64) ([]entity.Order, error)

	FindItemsByOrderID(ctx context.Context, orderID int64) ([]entity.OrderItem, error)

	Create(ctx context.Context, order *entity.Order) error
}
