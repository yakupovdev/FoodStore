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
	sellerRepo  repository.SellerRepository
	productRepo repository.ProductRepository
}

func NewSellerUsecase(sellerRepo repository.SellerRepository, productRepo repository.ProductRepository) (*SellerUsecase, error) {
	if sellerRepo == nil || productRepo == nil {
		return nil, domain.ErrDatabaseConnection
	}
	return &SellerUsecase{
		sellerRepo:  sellerRepo,
		productRepo: productRepo,
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
			ProductID:   o.ProductID,
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

func (uc *SellerUsecase) GetAllExistProducts(ctx context.Context) ([]dto.CategoryIDOutput, error) {
	result := make([]dto.CategoryIDOutput, 0)
	categories, err := uc.productRepo.GetCategories(ctx)
	if err != nil {
		return nil, err
	}
	for _, category := range categories {
		categoryID := dto.CategoryIDOutput{
			CategoryID:    category.ID,
			Name:          category.Name,
			SubCategories: make([]dto.SubCategoryIDOutput, 0),
		}
		subCategories, err := uc.productRepo.GetSubCategories(ctx, category.ID)
		if err != nil {
			return nil, err
		}
		for _, subCategory := range subCategories {
			subCategoryID := dto.SubCategoryIDOutput{
				SubCategoryID: subCategory.ID,
				Name:          subCategory.Name,
				Products:      make([]dto.ProductIDOutput, 0),
			}

			products, err := uc.productRepo.GetProductsBySubCategoryID(ctx, subCategory.ID)
			if err != nil {
				return nil, err
			}

			for _, product := range products {
				productID := dto.ProductIDOutput{
					ProductID:   product.ID,
					Name:        product.Name,
					Description: product.Description,
					Image:       product.Image,
				}
				subCategoryID.Products = append(subCategoryID.Products, productID)
			}

			categoryID.SubCategories = append(categoryID.SubCategories, subCategoryID)
		}

		result = append(result, categoryID)
	}

	return result, nil
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

func (uc *SellerUsecase) CreateOfferByExistProducts(ctx context.Context, input dto.CreateOfferByExistProductsInput) (*dto.CreateOfferByExistProductsOutput, error) {
	offerWithID, err := entity.NewOfferID(input.SellerID, input.CategoryID, input.SubCategoryID, input.ProductID, input.Price, input.Quantity)
	if err != nil {
		return nil, err
	}

	err = uc.sellerRepo.CreateOfferByExistProducts(ctx, offerWithID)
	if err != nil {
		return nil, err
	}

	res := &dto.CreateOfferByExistProductsOutput{
		Message:       "Created successfully",
		ProductID:     offerWithID.ProductID,
		CategoryID:    offerWithID.CategoryID,
		SubCategoryID: offerWithID.SubCategoryID,
		Price:         offerWithID.Price,
		Quantity:      offerWithID.Quantity,
	}

	return res, nil
}

func (uc *SellerUsecase) UpdateOffer(ctx context.Context, input dto.UpdateOfferInput) (*dto.UpdateOfferOutput, error) {
	offer, err := entity.NewSellerOffer(input.SellerID, input.ProductID, input.Price, input.Quantity)
	if err != nil {
		return nil, err
	}

	err = uc.sellerRepo.UpdateOffer(ctx, offer)
	if err != nil {
		return nil, err
	}

	return &dto.UpdateOfferOutput{
		Message:   "Offer updated successfully",
		ProductID: offer.ProductID,
		Price:     offer.Price,
		Quantity:  offer.Quantity,
	}, nil
}

func (uc *SellerUsecase) DeleteOffer(ctx context.Context, input dto.DeleteOfferInput) (*dto.DeleteOfferOutput, error) {
	offerPrimary, err := entity.NewOfferPrimary(input.SellerID, input.ProductID)
	if err != nil {
		return nil, err
	}
	err = uc.sellerRepo.DeleteOffer(ctx, offerPrimary)
	if err != nil {
		return nil, err
	}

	return &dto.DeleteOfferOutput{
		Message:   "Offer deleted successfully",
		ProductID: offerPrimary.ProductID,
	}, nil
}

func (uc *SellerUsecase) PurchaseSubscription(ctx context.Context, input dto.PurchaseSubscriptionInput) (*dto.PurchaseSubscriptionOutput, error) {
	params := &entity.PurchaseSubscriptionParams{
		ID:    input.ID,
		Price: 10000,
	}

	err := uc.sellerRepo.PurchaseSubscription(ctx, params)
	if err != nil {
		return nil, domain.ErrSubscriptionNotFound
	}

	return &dto.PurchaseSubscriptionOutput{
		Message: "Subscription purchased successfully",
	}, nil
}
