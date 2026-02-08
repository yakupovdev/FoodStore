package usecase

import (
	"context"
	"fmt"

	dto2 "github.com/yakupovdev/FoodStore/internal/delivery/http/dto"
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

func (uc *ClientUsecase) GetProfileByID(ctx context.Context, clientID int64, userType string) (*dto2.ClientProfileOutput, error) {
	client, err := uc.clientRepo.FindByID(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("get client profile: %w", err)
	}

	return &dto2.ClientProfileOutput{
		ID:       client.ID,
		Name:     client.Name,
		Email:    client.Email,
		UserType: userType,
		Balance:  client.Balance,
		Rating:   client.Rating,
	}, nil
}

func (uc *ClientUsecase) GetOrdersByClientID(ctx context.Context, clientID int64) ([]dto2.ClientOrderOutput, error) {
	orders, err := uc.orderRepo.FindByClientID(ctx, clientID)
	if err != nil {
		return nil, err
	}

	var result []dto2.ClientOrderOutput
	for _, order := range orders {
		items, err := uc.orderRepo.FindItemsByOrderID(ctx, order.ID)
		if err != nil {
			return nil, err
		}

		var itemDTOs []dto2.ClientOrderItemDTO
		for _, item := range items {
			itemDTOs = append(itemDTOs, dto2.ClientOrderItemDTO{
				OrderItemsID:    item.ID,
				OrderID:         item.OrderID,
				SellerID:        item.SellerID,
				ProductID:       item.ProductID,
				Quantity:        item.Quantity,
				PriceAtPurchase: item.PriceAtPurchase,
			})
		}

		result = append(result, dto2.ClientOrderOutput{
			OrderID:   order.ID,
			ClientID:  order.ClientID,
			Status:    order.Status,
			CreatedAt: order.CreatedAt,
			Items:     itemDTOs,
		})
	}

	return result, nil
}

func (uc *ClientUsecase) GetProducts(ctx context.Context) ([]dto2.CategoryOutput, error) {
	categories, err := uc.productRepo.GetCategories(ctx)
	if err != nil {
		return nil, err
	}

	var result []dto2.CategoryOutput
	for _, category := range categories {
		subCategories, err := uc.productRepo.GetSubCategories(ctx, category.ID)
		if err != nil {
			return nil, err
		}

		var subCatDTOs []dto2.SubCategoryOutput
		for _, subCat := range subCategories {
			products, err := uc.productRepo.GetProductsBySubCategoryID(ctx, subCat.ID)
			if err != nil {
				return nil, err
			}

			var productDTOs []dto2.ProductOutput
			for _, product := range products {
				offers, err := uc.sellerRepo.GetOffersByProductID(ctx, product.ID)
				if err != nil {
					return nil, err
				}

				var offerDTOs []dto2.SellerOfferOutput
				for _, offer := range offers {
					offerDTOs = append(offerDTOs, dto2.SellerOfferOutput{
						SellerID:   offer.SellerID,
						SellerName: offer.SellerName,
						Price:      offer.Price,
						Quantity:   offer.Quantity,
					})
				}

				productDTOs = append(productDTOs, dto2.ProductOutput{
					ProductID:   product.ID,
					Name:        product.Name,
					Description: product.Description,
					Image:       product.Image,
					Offers:      offerDTOs,
				})
			}

			subCatDTOs = append(subCatDTOs, dto2.SubCategoryOutput{
				Name:     subCat.Name,
				Products: productDTOs,
			})
		}

		result = append(result, dto2.CategoryOutput{
			Name:          category.Name,
			SubCategories: subCatDTOs,
		})
	}

	return result, nil
}
