package usecase

import (
	"context"
	"fmt"
	"log"

	"github.com/yakupovdev/FoodStore/internal/delivery/http/dto"
	"github.com/yakupovdev/FoodStore/internal/domain"
	"github.com/yakupovdev/FoodStore/internal/domain/entity"
	"github.com/yakupovdev/FoodStore/internal/domain/repository"
)

type ClientUsecase struct {
	clientRepo      repository.ClientRepository
	orderRepo       repository.OrderRepository
	productRepo     repository.ProductRepository
	sellerRepo      repository.SellerRepository
	transactionRepo repository.TransactionRepository
}

func NewClientUsecase(
	clientRepo repository.ClientRepository,
	orderRepo repository.OrderRepository,
	productRepo repository.ProductRepository,
	sellerRepo repository.SellerRepository,
	transactionRepo repository.TransactionRepository,
) (*ClientUsecase, error) {
	if clientRepo == nil || orderRepo == nil || productRepo == nil {
		return nil, domain.ErrDatabaseConnection
	}
	return &ClientUsecase{
		clientRepo:      clientRepo,
		orderRepo:       orderRepo,
		productRepo:     productRepo,
		sellerRepo:      sellerRepo,
		transactionRepo: transactionRepo,
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

func (uc *ClientUsecase) AddToCart(ctx context.Context, input dto.AddToCartInput) (*dto.AddToCartOutput, error) {

	offer, err := uc.sellerRepo.GetOffersBySellerIDAndProductID(ctx, input.SellerID, input.ProductID)
	if err != nil {
		return &dto.AddToCartOutput{}, fmt.Errorf("get offer price: %w", err)
	}
	if offer.Quantity < input.Quantity {
		return &dto.AddToCartOutput{}, fmt.Errorf("not enough stock available")
	}
	order := &entity.Order{
		ClientID: input.ClientID,
		Status:   "cart",
		Items: []entity.OrderItem{
			{
				SellerID:        input.SellerID,
				ProductID:       input.ProductID,
				Quantity:        input.Quantity,
				PriceAtPurchase: offer.Price,
			},
		},
	}

	if err := uc.clientRepo.AddToCart(ctx, *order); err != nil {
		return &dto.AddToCartOutput{}, fmt.Errorf("add to cart: %w", err)
	}
	return &dto.AddToCartOutput{
		Message: "item added to cart successfully",
	}, nil
}

func (uc *ClientUsecase) GetCartItems(ctx context.Context, clientID int64) ([]dto.CartItemOutput, error) {
	items, err := uc.clientRepo.GetCartItems(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("get cart items: %w", err)
	}

	var result []dto.CartItemOutput
	for _, item := range items {
		result = append(result, dto.CartItemOutput{
			CartItemsID:     item.ID,
			CartID:          item.OrderID,
			SellerID:        item.SellerID,
			ProductID:       item.ProductID,
			Quantity:        item.Quantity,
			PriceAtPurchase: item.PriceAtPurchase,
		})
	}

	return result, nil
}

func (uc *ClientUsecase) CreateOrder(ctx context.Context, input dto.CreateOrderInput) (*dto.CreateOrderOutput, error) {

	order := &entity.Order{
		ClientID: input.ClientID,
		Status:   "pending",
	}
	count := 0
	for _, item := range input.Items {

		offer, err := uc.sellerRepo.GetOffersBySellerIDAndProductID(ctx, item.SellerID, item.ProductID)
		if err != nil {
			log.Println(err, "error getting offer price in CreateOrder usecase")
			return nil, fmt.Errorf("get offer price: %w", err)
		}

		totalAmount := item.Quantity * offer.Price

		if err := uc.transactionRepo.ExecuteSellerTransaction(ctx, item.SellerID, totalAmount); err != nil {
			log.Println(err, "error executing seller transaction in CreateOrder usecase")
			return nil, fmt.Errorf("execute seller transaction: %w", err)
		}
		log.Println("seller transaction executed successfully for seller ID:", item.SellerID)
		order.Items = append(order.Items, entity.OrderItem{
			SellerID:        item.SellerID,
			ProductID:       item.ProductID,
			Quantity:        item.Quantity,
			PriceAtPurchase: offer.Price,
		})
		count++
	}

	if err := uc.orderRepo.Create(ctx, order); err != nil {
		log.Println(err, "error creating order in CreateOrder usecase")
		return nil, fmt.Errorf("create order: %w", err)
	}
	log.Println("order created successfully with ID:", order.ID)

	var totalCost int64
	for _, item := range input.Items {
		offer, err := uc.sellerRepo.GetOffersByProductID(ctx, item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("get offer price: %w", err)
		}
		priceAtPurchase := int64(0)
		for _, o := range offer {
			if o.SellerID == item.SellerID {
				priceAtPurchase = o.Price
				break
			}
		}
		totalCost += priceAtPurchase * int64(item.Quantity)
	}

	err := uc.transactionRepo.ExecuteOrderTransaction(ctx, input.ClientID, totalCost)
	if err != nil {
		if err == domain.ErrNotEnoughBalance {
			return nil, fmt.Errorf("not enough balance: %w", err)
		}
		log.Println(err, "error executing order transaction in CreateOrder usecase")
		return nil, fmt.Errorf("execute order transaction: %w", err)
	}

	return &dto.CreateOrderOutput{
		OrderID:   order.ID,
		ClientID:  order.ClientID,
		Status:    order.Status,
		CreatedAt: order.CreatedAt,
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

func (uc *ClientUsecase) UpdateBalance(ctx context.Context, input dto.UpdateBalanceInput) (*dto.BalanceUpdateOutput, error) {
	if err := uc.clientRepo.UpdateBalance(ctx, input.ClientID, input.Balance); err != nil {
		return &dto.BalanceUpdateOutput{}, fmt.Errorf("update balance: %w", err)
	}
	return &dto.BalanceUpdateOutput{
		Message: "balance updated successfully",
	}, nil
}

func (uc *ClientUsecase) AddAddress(ctx context.Context, input dto.AddAddressInput) (*dto.AddAddressOutput, error) {
	client := entity.Client{
		ID:       input.ClientID,
		Name:     "",
		Email:    "",
		UserType: "",
		Balance:  0,
		Rating:   0,
		Address:  input.Address,
	}
	if err := uc.clientRepo.AddAddress(ctx, client); err != nil {
		return &dto.AddAddressOutput{}, fmt.Errorf("add address: %w", err)
	}
	return &dto.AddAddressOutput{
		Message: "address added successfully",
	}, nil
}
