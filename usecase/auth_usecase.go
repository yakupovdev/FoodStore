package usecase

import (
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

func (au *AuthUsecase) LoginUser(email string, password string) (string, error) {
	hashHex := security.HashPassword(password)
	userID, err := au.repo.LoginUser(email, hashHex)

	if err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := security.GenerateToken(userID)

	if err != nil {
		return "", ErrTokenGeneration
	}

	return token, nil
}
