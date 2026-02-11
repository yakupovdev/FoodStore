package repository

import (
	"context"

	"github.com/yakupovdev/FoodStore/internal/domain/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) (int64, error)

	FindByEmailAndType(ctx context.Context, email, userType string) (*entity.User, error)

	ExistsByEmailAndType(ctx context.Context, email, userType string) (bool, error)

	FindIDByEmailAndType(ctx context.Context, email, userType string) (int64, error)

	UpdatePassword(ctx context.Context, userID int64, passwordHash string) error

	UpdateLastLogin(ctx context.Context, userID int64) error

	Delete(ctx context.Context, userID int64) error

	GetAllUsers(ctx context.Context) ([]entity.User, error)
}
