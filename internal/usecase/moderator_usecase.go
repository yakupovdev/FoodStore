package usecase

import (
	"context"

	"github.com/yakupovdev/FoodStore/internal/delivery/http/dto"
	"github.com/yakupovdev/FoodStore/internal/domain"
	"github.com/yakupovdev/FoodStore/internal/domain/entity"
	"github.com/yakupovdev/FoodStore/internal/domain/repository"
	"github.com/yakupovdev/FoodStore/internal/domain/service"
)

type ModeratorUsecase struct {
	moderatorRepo repository.ModeratorRepository
	productRepo   repository.ProductRepository
	sellerRepo    repository.SellerRepository
	emailSender   service.EmailSender
}

func NewModeratorUsecase(moderatorRepo repository.ModeratorRepository, productRepo repository.ProductRepository, sellerRepo repository.SellerRepository, emailSender service.EmailSender) (*ModeratorUsecase, error) {
	if moderatorRepo == nil || productRepo == nil || sellerRepo == nil || emailSender == nil {
		return nil, domain.ErrDatabaseConnection
	}
	return &ModeratorUsecase{
		moderatorRepo: moderatorRepo,
		productRepo:   productRepo,
		sellerRepo:    sellerRepo,
		emailSender:   emailSender,
	}, nil
}

func (uc *ModeratorUsecase) GetAllExistProducts(ctx context.Context) ([]dto.CategoryIDOutput, error) {
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

func (uc *ModeratorUsecase) GetModerationSellerOffers(ctx context.Context) ([]dto.OfferModerationOutput, error) {
	result := make([]dto.OfferModerationOutput, 0)

	offers, err := uc.moderatorRepo.GetModerationSellerOffers(ctx)
	if err != nil {
		return nil, err
	}

	for _, offer := range offers {
		o := dto.OfferModerationOutput{
			ProductID:       offer.ProductID,
			SellerID:        offer.SellerID,
			CategoryID:      offer.CategoryID,
			SubCategoryID:   offer.SubCategoryID,
			SellerName:      offer.SellerName,
			CategoryName:    offer.CategoryName,
			SubCategoryName: offer.SubCategoryName,
			ProductName:     offer.ProductName,
			Description:     offer.Description,
			Image:           offer.Image,
			Price:           offer.Price,
			Quantity:        offer.Quantity,
		}
		result = append(result, o)
	}
	return result, nil
}

func (uc *ModeratorUsecase) ApproveOffer(ctx context.Context, input dto.OfferModerationAnswerInput) (*dto.OfferModerationAnswerOutput, error) {
	offer, err := uc.moderatorRepo.GetSellerOfferByProductID(ctx, input.ProductID)
	if err != nil {
		return nil, err
	}

	product, err := entity.NewCreationProduct(offer.SubCategoryID, offer.ProductName, offer.Description, offer.Image)
	if err != nil {
		return nil, err
	}

	productID, err := uc.productRepo.CreateProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	offerWithID, err := entity.NewOfferWithID(offer.SellerID, offer.CategoryID, offer.SubCategoryID, productID, offer.Price, offer.Quantity)
	if err != nil {
		return nil, err
	}

	err = uc.sellerRepo.CreateOfferByExistProducts(ctx, offerWithID)
	if err != nil {
		return nil, err
	}

	msg := "Your offer " + offer.ProductName + "was successfully approved. Addition: " + input.Message

	err = uc.moderatorRepo.DeleteModerationSellerOffer(ctx, offer.ProductID)
	if err != nil {
		return nil, err
	}

	err = uc.emailSender.SendMessage(offer.SellerEmail, msg)
	if err != nil {
		return nil, err
	}

	return &dto.OfferModerationAnswerOutput{
		Message:   "Offer successfully approved",
		ProductID: input.ProductID,
	}, nil
}

func (uc *ModeratorUsecase) RejectOffer(ctx context.Context, input dto.OfferModerationAnswerInput) (*dto.OfferModerationAnswerOutput, error) {
	offer, err := uc.moderatorRepo.GetSellerOfferByProductID(ctx, input.ProductID)
	if err != nil {
		return nil, err
	}

	err = uc.moderatorRepo.DeleteModerationSellerOffer(ctx, offer.ProductID)
	if err != nil {
		return nil, err
	}

	msg := "Your offer " + offer.ProductName + "was rejected. Addition: " + input.Message

	err = uc.emailSender.SendMessage(offer.SellerEmail, msg)
	if err != nil {
		return nil, err
	}

	return &dto.OfferModerationAnswerOutput{
		Message:   "Offer successfully rejected",
		ProductID: input.ProductID,
	}, nil
}
