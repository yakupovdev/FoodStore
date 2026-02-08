package entity

type Category struct {
	ID            int64
	Name          string
	SubCategories []SubCategory
}

type SubCategory struct {
	ID       int64
	Name     string
	Products []Product
}

type Product struct {
	ID          int64
	Name        string
	Description string
	Image       string
}
