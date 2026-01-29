package usecase

import "github.com/yakupovdev/FoodStore/internal/repository"

type CheckTokenIsValidUsecase struct {
	repo *repository.Postgres
}

func NewCheckTokenIsValidUsecase(repo *repository.Postgres) (*CheckTokenIsValidUsecase, error) {
	if repo == nil {
		return nil, ErrDatabaseConnection
	}

	return &CheckTokenIsValidUsecase{
		repo: repo,
	}, nil
}

func (ctiv *CheckTokenIsValidUsecase) IsTokenValid(userID int64, token string) (bool, error) {
	return ctiv.repo.IsAccessTokenValid(userID, token)
}
