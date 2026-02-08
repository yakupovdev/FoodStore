package dto

import "time"

type ClientProfileOutput struct {
	ID       int64   `json:"id"`
	Name     string  `json:"name"`
	Email    string  `json:"email"`
	UserType string  `json:"user_type"`
	Balance  int64   `json:"balance"`
	Rating   float64 `json:"rating"`
}

type ClientOrderOutput struct {
	OrderID   int64                `json:"order_id"`
	ClientID  int64                `json:"client_id"`
	Status    string               `json:"status"`
	CreatedAt time.Time            `json:"created_at"`
	Items     []ClientOrderItemDTO `json:"items"`
}

type ClientOrderItemDTO struct {
	OrderItemsID    int64 `json:"order_items_id"`
	OrderID         int64 `json:"order_id"`
	SellerID        int64 `json:"seller_id"`
	ProductID       int64 `json:"product_id"`
	Quantity        int64 `json:"quantity"`
	PriceAtPurchase int64 `json:"price_at_purchase"`
}
