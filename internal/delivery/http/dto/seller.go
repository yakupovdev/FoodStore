package dto

type SellerProfileOutput struct {
	Name    string  `json:"name"`
	Email   string  `json:"email"`
	Type    string  `json:"type"`
	Balance int64   `json:"balance"`
	Rating  float32 `json:"rating"`
}

type SellerOfferOutput struct {
	SellerID   int64  `json:"seller_id"`
	SellerName string `json:"seller_name"`
	Price      int64  `json:"price"`
	Quantity   int64  `json:"quantity"`
}

type SellerOffersListOutput struct {
	Offers []SellerOfferItem `json:"offers"`
}

type SellerOfferItem struct {
	SellerID    int64  `json:"seller_id,omitempty"`
	ProductID   int64  `json:"product_id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Price       int64  `json:"price"`
	Quantity    int64  `json:"quantity"`
}

type CreateOfferInput struct {
	ProductName     string `json:"product_name"`
	Description     string `json:"description"`
	Image           string `json:"image"`
	Price           int64  `json:"price"`
	Quantity        int64  `json:"quantity"`
	CategoryName    string `json:"category_name"`
	SubCategoryName string `json:"sub_category_name"`
}

type PurchaseSubscriptionInput struct {
	ID int64 `json:"id"`
}

type CreateOfferOutput struct {
	Message         string `json:"message"`
	ProductName     string `json:"product_name"`
	Description     string `json:"description"`
	Image           string `json:"image"`
	Price           int64  `json:"price"`
	Quantity        int64  `json:"quantity"`
	CategoryName    string `json:"category_name"`
	SubCategoryName string `json:"sub_category_name"`
}

type CreateOfferByExistProductsInput struct {
	SellerID      int64 `json:"seller_id"`
	CategoryID    int64 `json:"category_id"`
	SubCategoryID int64 `json:"sub_category_id"`
	ProductID     int64 `json:"product_id"`
	Price         int64 `json:"price"`
	Quantity      int64 `json:"quantity"`
}

type CreateOfferByExistProductsOutput struct {
	Message       string `json:"message"`
	ProductID     int64  `json:"product_id"`
	CategoryID    int64  `json:"category_id"`
	SubCategoryID int64  `json:"sub_category_id"`
	Price         int64  `json:"price"`
	Quantity      int64  `json:"quantity"`
}

type UpdateOfferInput struct {
	SellerID  int64 `json:"seller_id"`
	ProductID int64 `json:"product_id"`
	Price     int64 `json:"price"`
	Quantity  int64 `json:"quantity"`
}

type UpdateOfferOutput struct {
	Message   string `json:"message"`
	ProductID int64  `json:"product_id"`
	Price     int64  `json:"price"`
	Quantity  int64  `json:"quantity"`
}

type DeleteOfferInput struct {
	SellerID  int64 `json:"seller_id"`
	ProductID int64 `json:"product_id"`
}

type DeleteOfferOutput struct {
	Message   string `json:"message"`
	ProductID int64  `json:"product_id"`
}

type PurchaseSubscriptionOutput struct {
	Message string `json:"message"`
}
