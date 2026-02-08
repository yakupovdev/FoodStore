package repository

import (
	"context"

	"github.com/yakupovdev/FoodStore/internal/domain/entity"
)

type ClientRepository interface {
	FindByID(ctx context.Context, clientID int64) (*entity.Client, error)
}
