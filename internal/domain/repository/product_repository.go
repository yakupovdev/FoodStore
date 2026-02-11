package repository

import (
	"context"

	"github.com/yakupovdev/FoodStore/internal/domain/entity"
)

type ProductRepository interface {
	GetCategories(ctx context.Context) ([]entity.Category, error)

	GetSubCategories(ctx context.Context, categoryID int64) ([]entity.SubCategory, error)

	GetProductsBySubCategoryID(ctx context.Context, subCategoryID int64) ([]entity.Product, error)

	GetProductByID(ctx context.Context, productID int64) (*entity.Product, error)

	GetProductsByPrioity(ctx context.Context) ([]entity.PriorityProduct, error)
}
