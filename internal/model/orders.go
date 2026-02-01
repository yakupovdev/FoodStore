package model

import "time"

type ClientOrders struct {
	OrderID   int64
	ClientID  int64
	Status    string
	CreatedAt time.Time
}

type ClientOrdersItems struct {
	OrderItemsId    int64
	OrderID         int64
	SellerID        int64
	ProductID       int64
	Quantity        int64
	PriceAtPurchase int64
}
