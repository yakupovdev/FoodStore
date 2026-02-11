package entity

import (
	"strings"

	"github.com/yakupovdev/FoodStore/internal/domain"
)

type ModerationOffer struct {
	ProductID       int64
	SellerID        int64
	CategoryID      int64
	SubCategoryID   int64
	SellerName      string
	SellerEmail     string
	CategoryName    string
	SubCategoryName string
	ProductName     string
	Description     string
	Image           string
	Price           int64
	Quantity        int64
}

func NewModerationOffer(sellerID, categoryID, subCategoryID int64, sellerName, sellerEmail, categoryName, subCategoryName, productName, description, image string, price, quantity int64) (*ModerationOffer, error) {
	productName = strings.TrimSpace(productName)
	description = strings.TrimSpace(description)
	image = strings.TrimSpace(image)

	if productName == "" {
		return nil, domain.ErrInvalidProductName
	}
	if description == "" {
		return nil, domain.ErrInvalidDescription
	}

	if price <= 0 {
		return nil, domain.ErrInvalidPrice
	}
	if quantity <= 0 {
		return nil, domain.ErrInvalidQuantity
	}
	return &ModerationOffer{
		SellerID:        sellerID,
		CategoryID:      categoryID,
		SubCategoryID:   subCategoryID,
		SellerName:      sellerName,
		SellerEmail:     sellerEmail,
		CategoryName:    categoryName,
		SubCategoryName: subCategoryName,
		ProductName:     productName,
		Description:     description,
		Image:           image,
		Price:           price,
		Quantity:        quantity,
	}, nil
}
