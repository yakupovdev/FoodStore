package model

type CategoryDTO struct {
	Name        string           `json:"name"`
	SubCategory []SubCategoryDTO `json:"sub_category"`
}

type SubCategoryDTO struct {
	Name     string       `json:"name"`
	Products []ProductDTO `json:"products"`
}

type ProductDTO struct {
	ProductID   int64      `json:"product_id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Image       string     `json:"image"`
	Offers      []OfferDTO `json:"sellers_offers"`
}
