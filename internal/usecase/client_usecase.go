package usecase

import (
	"github.com/yakupovdev/FoodStore/internal/model"
	"github.com/yakupovdev/FoodStore/internal/repository"
)

type ClientUsecase struct {
	repo  *repository.OrdersRepo
	repos *repository.SellerRepository
}

func NewClientUsecase(repo *repository.OrdersRepo, repos *repository.SellerRepository) (*ClientUsecase, error) {
	if repo == nil {
		return nil, ErrDatabaseConnection
	}

	return &ClientUsecase{
		repo:  repo,
		repos: repos,
	}, nil
}

func (ou *ClientUsecase) GetProfileByID(clientID int64) (model.Client, error) {
	profile, err := ou.repo.GetProfileByID(clientID)
	if err != nil {
		return model.Client{}, err
	}
	return profile, nil
}

func (ou *ClientUsecase) GetOrdersByClientID(clientID int64) ([]model.ClientOrdersDTO, error) {
	orders, err := ou.repo.GetOrdersByClientID(clientID)
	if err != nil {
		return nil, err
	}
	var detailedOrders []model.ClientOrdersDTO
	var detailedOrderItems []model.ClientOrdersItemDTO
	for _, order := range orders {
		items, err := ou.repo.GetOrderItemsByOrderID(order.OrderID)
		if err != nil {
			return nil, err
		}
		for _, item := range items {
			detailedItem := model.ClientOrdersItemDTO{
				OrderItemsId:    item.OrderItemsId,
				OrderID:         item.OrderID,
				SellerID:        item.SellerID,
				ProductID:       item.ProductID,
				Quantity:        item.Quantity,
				PriceAtPurchase: item.PriceAtPurchase,
			}
			detailedOrderItems = append(detailedOrderItems, detailedItem)
		}
		detailedOrder := model.ClientOrdersDTO{
			OrderID:   order.OrderID,
			ClientID:  order.ClientID,
			Status:    order.Status,
			CreatedAt: order.CreatedAt,
			Items:     detailedOrderItems,
		}
		detailedOrders = append(detailedOrders, detailedOrder)
	}
	return detailedOrders, nil
}

func (ou *ClientUsecase) GetOrderItemsByOrderID(orderID int64) ([]model.ClientOrdersItems, error) {
	items, err := ou.repo.GetOrderItemsByOrderID(orderID)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (ou *ClientUsecase) GetProducts() ([]model.CategoryDTO, error) {
	categories, err := ou.repo.GetCategories()
	if err != nil {
		return nil, err
	}

	var categoriesDTO []model.CategoryDTO
	for _, category := range categories {
		subCategories, err := ou.repo.GetSubCategoriesByCategoryID(category.ID)
		if err != nil {
			return nil, err
		}

		var subCategoriesDTO []model.SubCategoryDTO
		for _, subCategory := range subCategories {
			products, err := ou.repo.GetProductsBySubCategoryID(subCategory.ID)
			if err != nil {
				return nil, err
			}

			var productsDTO []model.ProductDTO
			for _, product := range products {
				sellers, err := ou.repos.GetSellerOffersByProductID(product.ID)
				if err != nil {
					return nil, err
				}
				var offersDTO []model.OfferDTO
				for _, seller := range sellers {
					offerDTO := model.OfferDTO{
						SellerID:   seller.SellerID,
						SellerName: seller.Name,
						Price:      seller.Price,
						Quantity:   seller.Quantity,
					}
					offersDTO = append(offersDTO, offerDTO)
				}
				productDTO := model.ProductDTO{
					ProductID:   product.ID,
					Name:        product.Name,
					Description: product.Description,
					Image:       product.Image,
					Offers:      offersDTO,
				}
				productsDTO = append(productsDTO, productDTO)
			}

			subCategoryDTO := model.SubCategoryDTO{
				Name:     subCategory.Name,
				Products: productsDTO,
			}
			subCategoriesDTO = append(subCategoriesDTO, subCategoryDTO)
		}

		categoryDTO := model.CategoryDTO{
			Name:        category.Name,
			SubCategory: subCategoriesDTO,
		}
		categoriesDTO = append(categoriesDTO, categoryDTO)
	}

	return categoriesDTO, nil
}
