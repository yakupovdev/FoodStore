package usecase

import (
	"context"
	"fmt"

	"github.com/yakupovdev/FoodStore/internal/delivery/http/dto"
	"github.com/yakupovdev/FoodStore/internal/domain"
	"github.com/yakupovdev/FoodStore/internal/domain/repository"
)

type ClientUsecase struct {
	clientRepo  repository.ClientRepository
	orderRepo   repository.OrderRepository
	productRepo repository.ProductRepository
	sellerRepo  repository.SellerRepository
}

func NewClientUsecase(
	clientRepo repository.ClientRepository,
	orderRepo repository.OrderRepository,
	productRepo repository.ProductRepository,
	sellerRepo repository.SellerRepository,
) (*ClientUsecase, error) {
	if clientRepo == nil || orderRepo == nil || productRepo == nil {
		return nil, domain.ErrDatabaseConnection
	}
	return &ClientUsecase{
		clientRepo:  clientRepo,
		orderRepo:   orderRepo,
		productRepo: productRepo,
		sellerRepo:  sellerRepo,
	}, nil
}

func (uc *ClientUsecase) GetProfileByID(ctx context.Context, clientID int64, userType string) (*dto.ClientProfileOutput, error) {
	client, err := uc.clientRepo.FindByID(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("get client profile: %w", err)
	}

	return &dto.ClientProfileOutput{
		ID:       client.ID,
		Name:     client.Name,
		Email:    client.Email,
		UserType: userType,
		Balance:  client.Balance,
		Rating:   client.Rating,
	}, nil
}

func (uc *ClientUsecase) GetOrdersByClientID(ctx context.Context, clientID int64) ([]dto.ClientOrderOutput, error) {
	orders, err := uc.orderRepo.FindByClientID(ctx, clientID)
	if err != nil {
		return nil, err
	}

	var result []dto.ClientOrderOutput
	for _, order := range orders {
		items, err := uc.orderRepo.FindItemsByOrderID(ctx, order.ID)
		if err != nil {
			return nil, err
		}

		var itemDTOs []dto.ClientOrderItemDTO
		for _, item := range items {
			itemDTOs = append(itemDTOs, dto.ClientOrderItemDTO{
				OrderItemsID:    item.ID,
				OrderID:         item.OrderID,
				SellerID:        item.SellerID,
				ProductID:       item.ProductID,
				Quantity:        item.Quantity,
				PriceAtPurchase: item.PriceAtPurchase,
			})
		}

		result = append(result, dto.ClientOrderOutput{
			OrderID:   order.ID,
			ClientID:  order.ClientID,
			Status:    order.Status,
			CreatedAt: order.CreatedAt,
			Items:     itemDTOs,
		})
	}

	return result, nil
}

func (uc *ClientUsecase) GetProducts(ctx context.Context) ([]dto.CategoryOutput, error) {
	categories, err := uc.productRepo.GetCategories(ctx)
	if err != nil {
		return nil, err
	}

	var result []dto.CategoryOutput
	for _, category := range categories {
		subCategories, err := uc.productRepo.GetSubCategories(ctx, category.ID)
		if err != nil {
			return nil, err
		}

		var subCatDTOs []dto.SubCategoryOutput
		for _, subCat := range subCategories {
			products, err := uc.productRepo.GetProductsBySubCategoryID(ctx, subCat.ID)
			if err != nil {
				return nil, err
			}

			var productDTOs []dto.ProductOutput
			for _, product := range products {
				offers, err := uc.sellerRepo.GetOffersByProductID(ctx, product.ID)
				if err != nil {
					return nil, err
				}

				var offerDTOs []dto.SellerOfferOutput
				for _, offer := range offers {
					offerDTOs = append(offerDTOs, dto.SellerOfferOutput{
						SellerID:   offer.SellerID,
						SellerName: offer.SellerName,
						Price:      offer.Price,
						Quantity:   offer.Quantity,
					})
				}

				productDTOs = append(productDTOs, dto.ProductOutput{
					ProductID:   product.ID,
					Name:        product.Name,
					Description: product.Description,
					Image:       product.Image,
					Offers:      offerDTOs,
				})
			}

			subCatDTOs = append(subCatDTOs, dto.SubCategoryOutput{
				Name:     subCat.Name,
				Products: productDTOs,
			})
		}

		result = append(result, dto.CategoryOutput{
			Name:          category.Name,
			SubCategories: subCatDTOs,
		})
	}

	return result, nil
}
