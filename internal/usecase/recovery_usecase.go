package usecase

import (
	"context"
	"fmt"

	"github.com/yakupovdev/FoodStore/internal/domain"
	"github.com/yakupovdev/FoodStore/internal/domain/entity"
	"github.com/yakupovdev/FoodStore/internal/domain/repository"
	"github.com/yakupovdev/FoodStore/internal/domain/service"
	"github.com/yakupovdev/FoodStore/internal/usecase/dto"
)

type RecoveryUsecase struct {
	userRepo     repository.UserRepository
	recoveryRepo repository.RecoveryCodeRepository
	codeHasher   service.CodeHasher
	tokenSvc     service.TokenService
	codeGen      service.CodeGenerator
	emailSender  service.EmailSender
}

func NewRecoveryUsecase(
	userRepo repository.UserRepository,
	recoveryRepo repository.RecoveryCodeRepository,
	codeHasher service.CodeHasher,
	tokenSvc service.TokenService,
	codeGen service.CodeGenerator,
	emailSender service.EmailSender,
) (*RecoveryUsecase, error) {
	if userRepo == nil || recoveryRepo == nil {
		return nil, domain.ErrDatabaseConnection
	}
	return &RecoveryUsecase{
		userRepo:     userRepo,
		recoveryRepo: recoveryRepo,
		codeHasher:   codeHasher,
		tokenSvc:     tokenSvc,
		codeGen:      codeGen,
		emailSender:  emailSender,
	}, nil
}

func (uc *RecoveryUsecase) SendCodeByEmail(ctx context.Context, input dto.SendCodeInput) error {
	userID, err := uc.userRepo.FindIDByEmailAndType(ctx, input.Email, input.UserType)
	if err != nil {
		return fmt.Errorf("%w: %w", domain.ErrUserNotFound, err)
	}

	code := uc.codeGen.GenerateRecoveryCode()

	if err := uc.emailSender.SendRecoveryCode(ctx, input.Email, code); err != nil {
		return fmt.Errorf("%w: %w", domain.ErrSMTPFailed, err)
	}

	codeHash := uc.codeHasher.Hash(code)

	if err := uc.recoveryRepo.Save(ctx, userID, input.Email, input.UserType, codeHash, nil); err != nil {
		return fmt.Errorf("save recovery code: %w", err)
	}

	return nil
}

func (uc *RecoveryUsecase) VerifyCode(ctx context.Context, input dto.VerifyCodeInput) (*dto.VerifyCodeOutput, error) {
	userID, err := uc.userRepo.FindIDByEmailAndType(ctx, input.Email, input.UserType)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", domain.ErrUserNotFound, err)
	}

	codeHash := uc.codeHasher.Hash(input.Code)

	isValid, err := uc.recoveryRepo.Verify(ctx, input.Email, input.UserType, codeHash)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", domain.ErrVerificationFailed, err)
	}

	if !isValid {
		return nil, domain.ErrCodeInvalid
	}

	recoveryToken, err := uc.tokenSvc.GenerateToken(userID, input.UserType, entity.RecoveryTokenType)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", domain.ErrTokenGeneration, err)
	}

	_ = uc.recoveryRepo.Delete(ctx, input.Email, input.UserType)

	return &dto.VerifyCodeOutput{
		RecoveryToken: recoveryToken,
		Message:       "Code verified successfully",
	}, nil
}

func (uc *RecoveryUsecase) ResetPassword(ctx context.Context, input dto.ResetPasswordInput) error {
	passwordHash, err := entity.HashPassword(input.NewPassword)
	if err != nil {
		return fmt.Errorf("%w: %w", domain.ErrPasswordHash, err)
	}

	if err := uc.userRepo.UpdatePassword(ctx, input.UserID, passwordHash); err != nil {
		return fmt.Errorf("%w: %w", domain.ErrUpdatePassword, err)
	}

	return nil
}
