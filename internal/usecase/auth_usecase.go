package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/yakupovdev/FoodStore/internal/delivery/http/dto"
	"github.com/yakupovdev/FoodStore/internal/domain"
	"github.com/yakupovdev/FoodStore/internal/domain/entity"
	"github.com/yakupovdev/FoodStore/internal/domain/repository"
	"github.com/yakupovdev/FoodStore/internal/domain/service"
)

type TokenValidator interface {
	IsTokenValid(ctx context.Context, userID int64, token string) (bool, error)
}

type AuthUsecase struct {
	userRepo  repository.UserRepository
	tokenRepo repository.TokenRepository
	tokenSvc  service.TokenService
}

func NewAuthUsecase(
	userRepo repository.UserRepository,
	tokenRepo repository.TokenRepository,
	tokenSvc service.TokenService,
) (*AuthUsecase, error) {
	if userRepo == nil || tokenRepo == nil {
		return nil, domain.ErrDatabaseConnection
	}
	return &AuthUsecase{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		tokenSvc:  tokenSvc,
	}, nil
}

func (uc *AuthUsecase) Register(ctx context.Context, input dto.RegisterInput) (*dto.RegisterOutput, error) {
	exists, err := uc.userRepo.ExistsByEmailAndType(ctx, input.Email, input.UserType)
	if err != nil {
		log.Printf("check user existence: %v", err)
		return nil, fmt.Errorf("check user existence: %w", err)
	}
	if exists {
		return nil, domain.ErrUserAlreadyExists
	}

	user, err := entity.NewUser(input.Email, input.Password, input.UserType, input.Name, input.Balance)
	if err != nil {
		log.Printf("create user entity: %v", err)
		return nil, fmt.Errorf("create user: %w", err)
	}

	if _, err = uc.userRepo.Create(ctx, user); err != nil {
		log.Printf("create user in db: %v", err)
		return nil, fmt.Errorf("persist user: %w", err)
	}

	return &dto.RegisterOutput{
		Message:  "User registered",
		Name:     input.Name,
		Email:    input.Email,
		UserType: input.UserType,
		Balance:  input.Balance,
	}, nil
}

func (uc *AuthUsecase) Login(ctx context.Context, input dto.LoginInput) (*dto.LoginOutput, error) {
	user, err := uc.userRepo.FindByEmailAndType(ctx, input.Email, input.UserType)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	if !user.CheckPassword(input.Password) {
		return nil, domain.ErrInvalidCredentials
	}

	accessToken, err := uc.tokenSvc.GenerateToken(user.ID, user.UserType, entity.AccessTokenType)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", domain.ErrTokenGeneration, err)
	}

	refreshToken, err := uc.tokenSvc.GenerateToken(user.ID, user.UserType, entity.RefreshTokenType)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", domain.ErrTokenGeneration, err)
	}

	if err := uc.tokenRepo.MoveToBlacklist(ctx, user.ID); err != nil {
		return nil, fmt.Errorf("%w: %w", domain.ErrTokenStorage, err)
	}

	expiredAt := time.Now().Add(1 * time.Hour)
	if err := uc.tokenRepo.SaveAccessToken(ctx, user.ID, accessToken, expiredAt); err != nil {
		return nil, fmt.Errorf("%w: %w", domain.ErrTokenStorage, err)
	}

	_ = uc.userRepo.UpdateLastLogin(ctx, user.ID)

	return &dto.LoginOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (uc *AuthUsecase) RefreshAccessToken(ctx context.Context, userID int64, userType string) (*dto.RefreshAccessTokenOutput, error) {
	accessToken, err := uc.tokenSvc.GenerateToken(userID, userType, entity.AccessTokenType)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", domain.ErrTokenGeneration, err)
	}

	if err := uc.tokenRepo.MoveToBlacklist(ctx, userID); err != nil {
		return nil, fmt.Errorf("%w: %w", domain.ErrTokenStorage, err)
	}

	if err := uc.tokenRepo.SaveAccessToken(ctx, userID, accessToken, time.Now().Add(1*time.Hour)); err != nil {
		return nil, fmt.Errorf("%w: %w", domain.ErrTokenStorage, err)
	}

	return &dto.RefreshAccessTokenOutput{
		AccessToken: accessToken,
	}, nil
}

func (uc *AuthUsecase) DeleteExpiredTokens(ctx context.Context) error {
	if err := uc.tokenRepo.DeleteExpiredTokens(ctx); err != nil {
		return fmt.Errorf("%w: %w", domain.ErrTokenCleanup, err)
	}
	return nil
}

func (uc *AuthUsecase) IsTokenValid(ctx context.Context, userID int64, token string) (bool, error) {
	return uc.tokenRepo.IsAccessTokenValid(ctx, userID, token)
}
