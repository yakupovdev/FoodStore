package usecase

import (
	"log"
	"time"

	"github.com/yakupovdev/FoodStore/internal/repository"
	"github.com/yakupovdev/FoodStore/internal/security"
)

type AuthUsecase struct {
	repo *repository.Postgres
}

func NewAuthUsecase(repo *repository.Postgres) (*AuthUsecase, error) {
	if repo == nil {
		return nil, ErrDatabaseConnection
	}

	return &AuthUsecase{
		repo: repo,
	}, nil
}

func (au *AuthUsecase) RegisterUser(email string, password string, userType string, balance int64) error {
	exist, err := au.repo.GetUserByEmail(email)

	if err != nil {
		return err
	}

	if exist {
		return ErrDuplicateEmail
	}

	hashHex := security.HashPassword(password)
	_, err = au.repo.RegisterUser(email, hashHex, userType, balance)

	if err != nil {
		return err
	}
	return nil
}

func (au *AuthUsecase) LoginUser(email string, password string) (string, string, error) {
	hashHex := security.HashPassword(password)
	userID, err := au.repo.LoginUser(email, hashHex)

	if err != nil {
		return "", "", ErrInvalidCredentials
	}

	accessToken, err := security.GenerateToken(userID,security.AccessToken)
	refreshToken, err := security.GenerateToken(userID,security.RefreshToken)
	expired_at := time.Now().Add(time.Hour)

	if err := au.repo.MoveFromWhiteListToBlackList(userID); err != nil {
		return "", "", ErrTokenStorage
	}
	if err := au.repo.SaveAccessToken(userID, accessToken, expired_at); err != nil {
		log.Println(err)
		return "", "", ErrTokenStorage
	}

	if err != nil {
		return "", "", ErrTokenGeneration
	}

	return accessToken, refreshToken, nil
}
