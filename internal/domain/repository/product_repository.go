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

	CreateProduct(ctx context.Context, product *entity.CreationProduct) (int64, error)

	GetParentID(ctx context.Context, productID int64) (int64, error)

	GetSubCategoryNameByID(ctx context.Context, categoryID int64) (string, error)

	GetCategoryNameByID(ctx context.Context, categoryID int64) (string, error)
  
	GetProductsByPrioity(ctx context.Context) ([]entity.PriorityProduct, error)
}
