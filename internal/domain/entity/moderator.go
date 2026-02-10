package entity

type ModeratorRequest struct {
	ID          int    `json:"id"`
	SellerID    int    `json:"seller_id"`
	SellerName  string `json:"seller_name"`
	ProductID   int    `json:"product_id"`
	ProductName string `json:"product_name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	Quantity    int    `json:"quantity"`
}
