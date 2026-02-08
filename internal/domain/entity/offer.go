package entity

import (
	"github.com/yakupovdev/FoodStore/internal/domain"
)

type Offer struct {
	SellerID    int64
	SellerName  string
	ProductName string
	Description string
	Image       string
	Price       int64
	Quantity    int64
}

type CreateOfferParams struct {
	SellerID        int64
	ProductName     string
	Description     string
	Image           string
	Price           int64
	Quantity        int64
	CategoryName    string
	SubCategoryName string
}

func NewCreateOfferParams(sellerID int64, productName, description, image string, price, quantity int64, categoryName, subCategoryName string) (*CreateOfferParams, error) {
	if productName == "" {
		return nil, domain.ErrEmptyProductName
	}
	if price <= 0 {
		return nil, domain.ErrInvalidPrice
	}
	if quantity < 0 {
		return nil, domain.ErrInvalidQuantity
	}
	if categoryName == "" {
		return nil, domain.ErrEmptyCategoryName
	}
	if subCategoryName == "" {
		return nil, domain.ErrEmptySubCatName
	}

	return &CreateOfferParams{
		SellerID:        sellerID,
		ProductName:     productName,
		Description:     description,
		Image:           image,
		Price:           price,
		Quantity:        quantity,
		CategoryName:    categoryName,
		SubCategoryName: subCategoryName,
	}, nil
}
