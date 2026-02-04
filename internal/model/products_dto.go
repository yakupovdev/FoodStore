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
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
}
