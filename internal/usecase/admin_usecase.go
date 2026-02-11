package usecase

import (
	"context"
	"fmt"
	"log"

	"github.com/yakupovdev/FoodStore/internal/delivery/http/dto"
	"github.com/yakupovdev/FoodStore/internal/domain/entity"
	"github.com/yakupovdev/FoodStore/internal/domain/repository"
	"github.com/yakupovdev/FoodStore/internal/domain/service"
)

type AdminUsecase struct {
	userRepo   repository.UserRepository
	clientRepo repository.ClientRepository
	adminRepo  repository.AdminRepository
	logsRepo   repository.LogsRepository
	ca         service.CheckerAdminKey
}

func NewAdminUsecase(
	userRepo repository.UserRepository,
	clientRepo repository.ClientRepository,
	adminRepo repository.AdminRepository,
	logsRepo repository.LogsRepository,
	ca service.CheckerAdminKey,
) (*AdminUsecase, error) {
	if userRepo == nil || clientRepo == nil || adminRepo == nil || logsRepo == nil || ca == nil {
		return nil, fmt.Errorf("database connection error")
	}
	return &AdminUsecase{
		userRepo:   userRepo,
		clientRepo: clientRepo,
		adminRepo:  adminRepo,
		logsRepo:   logsRepo,
		ca:         ca,
	}, nil
}

func (a *AdminUsecase) DeleteUser(ctx context.Context, params dto.DeleteUserInput) (dto.DeleteUserOutput, error) {
	log.Println(params.UserID)
	if err := a.userRepo.Delete(ctx, params.UserID); err != nil {
		return dto.DeleteUserOutput{}, err
	}
	output := dto.DeleteUserOutput{
		Message: "User deleted successfully",
	}
	return output, nil
}

func (a *AdminUsecase) UpdateBalance(ctx context.Context, input dto.UpdateBalanceInput) (*dto.BalanceUpdateOutput, error) {
	if err := a.adminRepo.UpdateBalance(ctx, input.UserID, input.Balance); err != nil {
		return &dto.BalanceUpdateOutput{}, fmt.Errorf("update balance: %w", err)
	}
	return &dto.BalanceUpdateOutput{
		Message: "balance updated successfully",
	}, nil
}

func (a *AdminUsecase) GetAllLogTransactions(ctx context.Context) ([]dto.GetLogsHistoryOutput, error) {
	logs, err := a.logsRepo.GetLogsHistory(ctx)
	if err != nil {
		return nil, fmt.Errorf("get logs history: %w", err)
	}

	var output []dto.GetLogsHistoryOutput
	for _, log := range logs {
		output = append(output, dto.GetLogsHistoryOutput{
			LogID:            log.ID,
			ClientID:         log.ClientID,
			SellerID:         log.SellerID,
			TotalAmount:      log.TotalAmount,
			CommissionAmount: log.CommissionAmount,
			CreatedAt:        log.CreatedAt,
		})
	}
	return output, nil
}

func (a *AdminUsecase) GetAllUsers(ctx context.Context) ([]dto.GetAllUsersOutput, error) {
	users, err := a.userRepo.GetAllUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all users: %w", err)
	}

	var output []dto.GetAllUsersOutput
	for _, user := range users {
		output = append(output, dto.GetAllUsersOutput{
			UserID:        user.ID,
			Email:         user.Email,
			UserType:      user.UserType,
			Balance:       user.Balance,
			CreatedAt:     user.CreatedAt,
			LastEnteredAt: user.LastEnter,
		})
	}
	return output, nil
}

func (a *AdminUsecase) CreateAdmin(ctx context.Context, input dto.CreateAdminInput) (*dto.CreateAdminOutput, error) {
	ok := a.ca.CheckAdminKey(input.SecretKey)
	if !ok {
		return nil, fmt.Errorf("invalid admin key")
	}
	user, err := entity.NewUser(input.Email, input.Password, "admin", "admin", 0)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	userID, err := a.userRepo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("create admin: %w", err)
	}

	return &dto.CreateAdminOutput{
		UserID:  userID,
		Message: "Admin created successfully",
	}, nil
}

func (a *AdminUsecase) DeleteExpiredSubscription(ctx context.Context) error {
	ids, err := a.adminRepo.GetExpiringSubscriptions(ctx)
	if err != nil {
		return fmt.Errorf("get expiring subscriptions: %w", err)
	}

	for _, id := range ids {
		if err := a.adminRepo.CancelSubscription(ctx, id); err != nil {
			return fmt.Errorf("cancel subscription: %w", err)
		}
		if err := a.adminRepo.SetPriorityToFalse(ctx, id); err != nil {
			return fmt.Errorf("set priority to false: %w", err)
		}
	}

	return nil
}
