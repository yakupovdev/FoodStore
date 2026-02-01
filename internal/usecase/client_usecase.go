package usecase

import (
	"github.com/yakupovdev/FoodStore/internal/model"
	"github.com/yakupovdev/FoodStore/internal/repository"
)

type ClientUsecase struct {
	repo *repository.OrdersRepo
}

func NewClientUsecase(repo *repository.OrdersRepo) (*ClientUsecase, error) {
	if repo == nil {
		return nil, ErrDatabaseConnection
	}

	return &ClientUsecase{
		repo: repo,
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
