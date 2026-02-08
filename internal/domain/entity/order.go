package entity

import "time"

type Order struct {
	ID        int64
	ClientID  int64
	Status    string
	CreatedAt time.Time
	Items     []OrderItem
}

type OrderItem struct {
	ID              int64
	OrderID         int64
	SellerID        int64
	ProductID       int64
	Quantity        int64
	PriceAtPurchase int64
}
