package impl

import (
	"context"
	"errors"
	"fmt"
	"log"

	pg "github.com/jackc/pgx/v5"
	"github.com/yakupovdev/FoodStore/internal/domain"
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

func (r *ProductRepo) GetProductByID(ctx context.Context, productID int64) (*entity.Product, error) {
	stmt := `SELECT product_id, name, description, img FROM products WHERE product_id=$1`
	row := r.conn.QueryRow(ctx, stmt, productID)

	var product entity.Product
	if err := row.Scan(&product.ID, &product.Name, &product.Description, &product.Image); err != nil {
		log.Println(err)
		if errors.Is(err, pg.ErrNoRows) {
			return nil, domain.ErrNoProducts
		}
		return nil, fmt.Errorf("get product by id: %w", err)
	}

	return &product, nil
}

func (r *ProductRepo) CreateProduct(ctx context.Context, product *entity.CreationProduct) (int64, error) {
	stmt := `
		INSERT INTO products (name, description, img, categories_id)
		VALUES ($1, $2, $3, $4)
		RETURNING product_id
	`

	var productID int64
	err := r.conn.
		QueryRow(ctx, stmt,
			product.Name,
			product.Description,
			product.Image,
			product.CategoryID,
		).
		Scan(&productID)

	if err != nil {
		log.Println(err)
		return 0, fmt.Errorf("create product: %w", err)
	}

	return productID, nil
}

func (r *ProductRepo) GetCategoryNameByID(ctx context.Context, categoryID int64) (string, error) {
	stmt := `SELECT name FROM categories WHERE categories_id=$1 AND parent_id IS NULL`

	var name string
	err := r.conn.QueryRow(ctx, stmt, categoryID).Scan(&name)
	if err != nil {
		log.Println(err)
		if errors.Is(err, pg.ErrNoRows) {
			return "", domain.ErrCategoryNotFound
		}
	}

	return name, nil
}

func (r *ProductRepo) GetSubCategoryNameByID(ctx context.Context, categoryID int64) (string, error) {
	stmt := `SELECT name FROM categories WHERE categories_id=$1 AND parent_id IS NOT NULL`

	var name string
	err := r.conn.QueryRow(ctx, stmt, categoryID).Scan(&name)
	if err != nil {
		log.Println(err)
		if errors.Is(err, pg.ErrNoRows) {
			return "", domain.ErrSubCategoryNotFound
		}
	}

	return name, nil
}

func (r *ProductRepo) GetParentID(ctx context.Context, categoryID int64) (int64, error) {
	stmt := `SELECT parent_id FROM categories WHERE categories_id=$1 AND parent_id IS NOT NULL`

	var parentID int64
	err := r.conn.QueryRow(ctx, stmt, categoryID).Scan(&parentID)
	if err != nil {
		log.Println(err)
		if errors.Is(err, pg.ErrNoRows) {
			return 0, domain.ErrSubCategoryNotFound
		}
	}

	return parentID, nil
}

func (r *ProductRepo) GetProductsByPrioity(ctx context.Context) ([]entity.PriorityProduct, error) {
	stmt := `SELECT 
    so.seller_id,
    s.name AS seller_name,
    s.priority,
    so.product_id,
    p.name AS product_name,
    p.description,
    so.price,
    so.quantity,
    p.img
	FROM seller_offers so
	INNER JOIN sellers s 
    ON so.seller_id = s.seller_id
	INNER JOIN products p 
    ON so.product_id = p.product_id
	WHERE s.priority = 1 and so.quantity <> 0;`

	rows, err := r.conn.Query(ctx, stmt)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("get products by priority: %w", err)
	}
	defer rows.Close()

	var priorityProducts []entity.PriorityProduct
	for rows.Next() {
		var pp entity.PriorityProduct
		if err := rows.Scan(&pp.SellerID, &pp.SellerName, &pp.Priority, &pp.ID, &pp.ProductName, &pp.Description, &pp.Price, &pp.Quantity, &pp.Img); err != nil {
			log.Println(err)
			return nil, fmt.Errorf("get products by priority: %w", err)
		}
		priorityProducts = append(priorityProducts, pp)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, fmt.Errorf("get products by priority: %w", err)
	}
	if len(priorityProducts) == 0 {
		return nil, domain.ErrNoProducts
	}

	return priorityProducts, nil
}
