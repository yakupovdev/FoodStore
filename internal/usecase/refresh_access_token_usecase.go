package usecase

import (
	"time"

	"github.com/yakupovdev/FoodStore/internal/repository"
	"github.com/yakupovdev/FoodStore/internal/security"
)

type RefreshAccessTokenUsecase struct {
	repo *repository.Postgres
}

func NewRefreshAccessTokenUsecase(repo *repository.Postgres) (*RefreshAccessTokenUsecase, error) {
	if repo == nil {
		return nil, ErrDatabaseConnection
	}

	return &RefreshAccessTokenUsecase{
		repo: repo,
	}, nil
}

func (rtu *RefreshAccessTokenUsecase) RefreshAccessToken(userID int64) (string, error) {
	accessToken, err := security.GenerateToken(userID,security.AccessToken)
	if err != nil {
		return "", ErrTokenGeneration
	}
	if err := rtu.repo.MoveFromWhiteListToBlackList(userID); err != nil {
		return "", ErrTokenStorage
	}
	if err := rtu.repo.SaveAccessToken(userID, accessToken, time.Now().Add(24*time.Hour)); err != nil {
		return "", ErrTokenStorage
	}

	return accessToken, nil
}

func (rtu *RefreshAccessTokenUsecase) DeleteExpiredTokens() error {
	if err := rtu.repo.DeleteExpiredAccessTokens(); err != nil {
		return ErrTokenCleanup
	}
	return nil
}
