package postgres

import (
	"context"
	"fmt"
	"log"

	pg "github.com/jackc/pgx/v5"
	"github.com/yakupovdev/FoodStore/internal/domain/entity"
)

type ProductRepo struct {
	conn *pg.Conn
}

func NewProductRepo(conn *pg.Conn) *ProductRepo {
	return &ProductRepo{conn: conn}
}

func (r *ProductRepo) GetCategories(ctx context.Context) ([]entity.Category, error) {
	stmt := `SELECT categories_id, name FROM categories WHERE parent_id IS NULL`
	rows, err := r.conn.Query(ctx, stmt)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("get categories: %w", err)
	}
	defer rows.Close()

	var categories []entity.Category
	for rows.Next() {
		var category entity.Category
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			log.Println(err)
			return nil, fmt.Errorf("get categories: %w", err)
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, fmt.Errorf("get categories: %w", err)
	}

	return categories, nil
}

func (r *ProductRepo) GetSubCategories(ctx context.Context, categoryID int64) ([]entity.SubCategory, error) {
	stmt := `SELECT categories_id, name FROM categories WHERE parent_id=$1`
	rows, err := r.conn.Query(ctx, stmt, categoryID)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("get subcategories: %w", err)
	}
	defer rows.Close()

	var subcategories []entity.SubCategory
	for rows.Next() {
		var subcategory entity.SubCategory
		if err := rows.Scan(&subcategory.ID, &subcategory.Name); err != nil {
			log.Println(err)
			return nil, fmt.Errorf("get subcategories: %w", err)
		}
		subcategories = append(subcategories, subcategory)
	}

	return subcategories, nil
}

func (r *ProductRepo) GetProductsBySubCategoryID(ctx context.Context, subCategoryID int64) ([]entity.Product, error) {
	stmt := `SELECT product_id, name, description, img FROM products WHERE categories_id=$1`
	rows, err := r.conn.Query(ctx, stmt, subCategoryID)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("get products: %w", err)
	}
	defer rows.Close()

	var products []entity.Product
	for rows.Next() {
		var product entity.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Image); err != nil {
			log.Println(err)
			return nil, fmt.Errorf("get products: %w", err)
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, fmt.Errorf("get products: %w", err)
	}

	return products, nil
}
