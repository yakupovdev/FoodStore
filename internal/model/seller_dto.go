package model

type SellerProfileResponse struct {
	Name    string  `json:"name"`
	Email   string  `json:"email"`
	Type    string  `json:"type"`
	Balance int64   `json:"balance"`
	Rating  float32 `json:"rating"`
}

type OfferDTO struct {
	SellerID   int64  `json:"seller_id"`
	SellerName string `json:"seller_name"`
	Price      int64  `json:"price"`
	Quantity   int64  `json:"quantity"`
}

type SellerOffersResponse struct {
	Offers []Offer `json:"offers"`
}

type CreateSellerOfferRequest struct {
	ProductName     string `json:"product_name"`
	Description     string `json:"description"`
	Image           string `json:"image"`
	Price           int64  `json:"price"`
	Quantity        int64  `json:"quantity"`
	CategoryName    string `json:"category_name"`
	SubCategoryName string `json:"sub_category_name"`
}
type CreateSellerOfferResponse struct {
	Message         string `json:"message"`
	ProductName     string `json:"product_name"`
	Description     string `json:"description"`
	Image           string `json:"image"`
	Price           int64  `json:"price"`
	Quantity        int64  `json:"quantity"`
	CategoryName    string `json:"category_name"`
	SubCategoryName string `json:"sub_category_name"`
}
