package repository

import (
	"context"

	"github.com/yakupovdev/FoodStore/internal/domain/entity"
)

type ClientRepository interface {
	FindByID(ctx context.Context, clientID int64) (*entity.Client, error)

	GetBalance(ctx context.Context, clientID int64) (int64, error)

	UpdateBalance(ctx context.Context, clientID int64, newBalance int64) error

	AddAddress(ctx context.Context, input entity.Client) error
}
