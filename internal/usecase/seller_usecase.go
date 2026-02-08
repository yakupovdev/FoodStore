package usecase

import (
	"context"
	"fmt"

	"github.com/yakupovdev/FoodStore/internal/delivery/http/dto"
	"github.com/yakupovdev/FoodStore/internal/domain"
	"github.com/yakupovdev/FoodStore/internal/domain/entity"
	"github.com/yakupovdev/FoodStore/internal/domain/repository"
)

type SellerUsecase struct {
	sellerRepo repository.SellerRepository
}

func NewSellerUsecase(sellerRepo repository.SellerRepository) (*SellerUsecase, error) {
	if sellerRepo == nil {
		return nil, domain.ErrDatabaseConnection
	}
	return &SellerUsecase{
		sellerRepo: sellerRepo,
	}, nil
}

func (uc *SellerUsecase) GetProfileByID(ctx context.Context, userID int64) (*dto.SellerProfileOutput, error) {
	seller, err := uc.sellerRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get seller profile: %w", err)
	}

	return &dto.SellerProfileOutput{
		Name:    seller.Name,
		Email:   seller.Email,
		Type:    "seller",
		Balance: seller.Balance,
		Rating:  seller.Rating,
	}, nil
}

func (uc *SellerUsecase) GetOffersByID(ctx context.Context, userID int64) (*dto.SellerOffersListOutput, error) {
	offers, err := uc.sellerRepo.GetOffersBySellerID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get seller offers: %w", err)
	}

	var items []dto.SellerOfferItem
	for _, o := range offers {
		items = append(items, dto.SellerOfferItem{
			Name:        o.ProductName,
			Description: o.Description,
			Image:       o.Image,
			Price:       o.Price,
			Quantity:    o.Quantity,
		})
	}

	return &dto.SellerOffersListOutput{
		Offers: items,
	}, nil
}

func (uc *SellerUsecase) CreateOffer(ctx context.Context, input dto.CreateOfferInput, sellerID int64) (*dto.CreateOfferOutput, error) {
	params, err := entity.NewCreateOfferParams(
		sellerID,
		input.ProductName,
		input.Description,
		input.Image,
		input.Price,
		input.Quantity,
		input.CategoryName,
		input.SubCategoryName,
	)
	if err != nil {
		return nil, err
	}

	if err := uc.sellerRepo.CreateOffer(ctx, params); err != nil {
		return nil, fmt.Errorf("create offer: %w", err)
	}

	return &dto.CreateOfferOutput{
		Message:         "Created successfully",
		ProductName:     input.ProductName,
		Description:     input.Description,
		Image:           input.Image,
		Price:           input.Price,
		Quantity:        input.Quantity,
		CategoryName:    input.CategoryName,
		SubCategoryName: input.SubCategoryName,
	}, nil
}
