package service

import (
	"context"

	"github.com/yakupovdev/FoodStore/internal/domain/entity"
)

type TokenService interface {
	GenerateToken(userID int64, userType string, tokenType entity.TokenType) (string, error)

	ParseToken(tokenStr string, tokenType entity.TokenType) (*entity.TokenClaims, error)
}

type CodeHasher interface {
	Hash(data string) string
}

type CodeGenerator interface {
	GenerateRecoveryCode() string
}

type EmailSender interface {
	SendRecoveryCode(ctx context.Context, email, code string) error
}
