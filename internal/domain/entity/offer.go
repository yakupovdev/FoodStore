package entity

import (
	"github.com/yakupovdev/FoodStore/internal/domain"
)

type Offer struct {
	SellerID    int64
	ProductID   int64
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

type OfferWithID struct {
	SellerID      int64
	CategoryID    int64
	SubCategoryID int64
	ProductID     int64
	Price         int64
	Quantity      int64
}

func NewOfferID(sellerID, categoryID, subCategoryID, productID, price, quantity int64) (*OfferWithID, error) {
	if categoryID <= 0 {
		return nil, domain.ErrCategoryID
	}
	if subCategoryID <= 0 {
		return nil, domain.ErrSubCategoryID
	}
	if productID <= 0 {
		return nil, domain.ErrProductID
	}
	if sellerID <= 0 {
		return nil, domain.ErrSellerID
	}
	if quantity <= 0 {
		return nil, domain.ErrInvalidQuantity
	}
	if price <= 0 {
		return nil, domain.ErrInvalidPrice
	}

	return &OfferWithID{
		SellerID:      sellerID,
		CategoryID:    categoryID,
		SubCategoryID: subCategoryID,
		ProductID:     productID,
		Price:         price,
		Quantity:      quantity,
	}, nil
}

type SellerOffer struct {
	SellerID  int64
	ProductID int64
	Price     int64
	Quantity  int64
}

func NewSellerOffer(sellerID int64, productID int64, price int64, quantity int64) (*SellerOffer, error) {
	if productID <= 0 {
		return nil, domain.ErrProductID
	}
	if price <= 0 {
		return nil, domain.ErrInvalidPrice
	}
	if sellerID <= 0 {
		return nil, domain.ErrSellerID
	}
	if quantity <= 0 {
		return nil, domain.ErrInvalidQuantity
	}

	return &SellerOffer{
		SellerID:  sellerID,
		ProductID: productID,
		Price:     price,
		Quantity:  quantity,
	}, nil
}

type OfferPrimary struct {
	SellerID  int64
	ProductID int64
}

func NewOfferPrimary(sellerID int64, productID int64) (*OfferPrimary, error) {
	if productID <= 0 {
		return nil, domain.ErrProductID
	}
	if sellerID <= 0 {
		return nil, domain.ErrSellerID
	}

	return &OfferPrimary{
		SellerID:  sellerID,
		ProductID: productID,
	}, nil
}
