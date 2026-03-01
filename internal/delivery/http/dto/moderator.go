package dto

type OfferModerationOutput struct {
	ProductID       int64  `json:"product_id"`
	SellerID        int64  `json:"seller_id"`
	CategoryID      int64  `json:"category_id"`
	SubCategoryID   int64  `json:"sub_category_id"`
	SellerName      string `json:"seller_name"`
	CategoryName    string `json:"category_name"`
	SubCategoryName string `json:"sub_category_name"`
	ProductName     string `json:"product_name"`
	Description     string `json:"description"`
	Image           string `json:"image"`
	Price           int64  `json:"price"`
	Quantity        int64  `json:"quantity"`
}

type OfferModerationAnswerInput struct {
	ProductID int64  `json:"product_id"`
	Message   string `json:"message"`
}

type OfferModerationAnswerOutput struct {
	Message   string `json:"message"`
	ProductID int64  `json:"product_id"`
}
