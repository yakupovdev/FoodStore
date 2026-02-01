package model

type Seller struct {
	Name    string
	Email   string
	Type    string
	Balance int64
	Rating  float32
}

type SellerOffers struct {
	Offers []Offer
}

type Offer struct {
	Name        string
	Description string
	Image       string
	Price       int64
	Quantity    int64
}

type CreateOfferParams struct {
	SellerID    int64
	ProductName string
	Description string
	Image       string
	Price       int64
	Quantity    int64

	CategoryName    string
	SubCategoryName string
}
