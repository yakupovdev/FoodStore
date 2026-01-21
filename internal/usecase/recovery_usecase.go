package usecase

import (
	"github.com/yakupovdev/FoodStore/internal/repository"
	"github.com/yakupovdev/FoodStore/internal/security"
)

type RecoveryUsecase struct {
	repo *repository.Postgres
}

func NewRecoveryUsecase(repo *repository.Postgres) (*RecoveryUsecase, error) {
	if repo == nil {
		return nil, ErrDatabaseConnection
	}

	return &RecoveryUsecase{
		repo: repo,
	}, nil
}

func (ru *RecoveryUsecase) ResetUserPassword(userID int64, newPassword string) error {
	hashHex := security.HashPassword(newPassword)
	err := ru.repo.UpdateUserPassword(userID, hashHex)

	if err != nil {
		return ErrUpdatePassword
	}

	return nil
}
