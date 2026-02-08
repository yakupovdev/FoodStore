package dto

type CategoryOutput struct {
	Name          string              `json:"name"`
	SubCategories []SubCategoryOutput `json:"sub_category"`
}

type SubCategoryOutput struct {
	Name     string          `json:"name"`
	Products []ProductOutput `json:"products"`
}

type ProductOutput struct {
	ProductID   int64               `json:"product_id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Image       string              `json:"image"`
	Offers      []SellerOfferOutput `json:"sellers_offers"`
}

type CategoryIDOutput struct {
	CategoryID    int64                 `json:"category_id"`
	Name          string                `json:"name"`
	SubCategories []SubCategoryIDOutput `json:"sub_category"`
}

type SubCategoryIDOutput struct {
	SubCategoryID int64             `json:"sub_category_id"`
	Name          string            `json:"name"`
	Products      []ProductIDOutput `json:"products"`
}

type ProductIDOutput struct {
	ProductID   int64  `json:"product_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
}
