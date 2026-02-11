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

type PriorityProduct struct {
	SellerID    int64
	SellerName  string
	Priority    int64
	ID          int64
	ProductName string
	Description string
	Price       int64
	Quantity    int64
	Img         string
}
