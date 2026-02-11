package entity

import (
	"strings"

	"github.com/yakupovdev/FoodStore/internal/domain"
)

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

type CreationProduct struct {
	CategoryID  int64
	Name        string
	Description string
	Image       string
}

func NewCreationProduct(categoryID int64, name string, description string, image string) (*CreationProduct, error) {
	name = strings.TrimSpace(name)
	description = strings.TrimSpace(description)
	image = strings.TrimSpace(image)

	if categoryID <= 0 {
		return nil, domain.ErrCategoryID
	}

	return &CreationProduct{
		CategoryID:  categoryID,
		Name:        name,
		Description: description,
		Image:       image,
	}, nil
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
