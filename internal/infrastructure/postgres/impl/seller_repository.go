package impl

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	pg "github.com/jackc/pgx/v5"
	"github.com/yakupovdev/FoodStore/internal/domain"
	"github.com/yakupovdev/FoodStore/internal/domain/entity"
)

type SellerRepo struct {
	conn *pg.Conn
}

func NewSellerRepo(conn *pg.Conn) *SellerRepo {
	return &SellerRepo{conn: conn}
}

func (r *SellerRepo) FindByUserID(ctx context.Context, userID int64) (*entity.Seller, error) {
	stmtBase := `
        SELECT u.email, u.balance, s.name, s.rating 
        FROM users u
        JOIN sellers s ON u.userid = s.seller_id
        WHERE u.userid = $1`

	var seller entity.Seller
	err := r.conn.QueryRow(ctx, stmtBase, userID).Scan(
		&seller.Email, &seller.Balance, &seller.Name, &seller.Rating,
	)
	if err != nil {
		return nil, fmt.Errorf("find seller: %w", err)
	}

	seller.ID = userID
	return &seller, nil
}

func (r *SellerRepo) GetOffersBySellerID(ctx context.Context, sellerID int64) ([]entity.Offer, error) {
	stmtOffers := `
        SELECT p.name, p.description, p.img, so.price, so.quantity
        FROM seller_offers so
        JOIN products p ON so.product_id = p.product_id
        WHERE so.seller_id = $1`

	rows, err := r.conn.Query(ctx, stmtOffers, sellerID)
	if err != nil {
		return make([]entity.Offer, 0), fmt.Errorf("query offers: %w", err)
	}
	defer rows.Close()

	var offers []entity.Offer
	for rows.Next() {
		var o entity.Offer
		if err := rows.Scan(&o.ProductName, &o.Description, &o.Image, &o.Price, &o.Quantity); err != nil {
			return make([]entity.Offer, 0), err
		}
		offers = append(offers, o)
	}

	return offers, nil
}

func (r *SellerRepo) GetOffersByProductID(ctx context.Context, productID int64) ([]entity.Offer, error) {
	stmtOffers := `
       SELECT s.seller_id, s.name, p.description, p.img, so.price, so.quantity
		FROM seller_offers so
		JOIN products p ON so.product_id = p.product_id
		JOIN sellers s ON s.seller_id = so.seller_id
		WHERE p.product_id = $1`

	rows, err := r.conn.Query(ctx, stmtOffers, productID)
	if err != nil {
		return make([]entity.Offer, 0), fmt.Errorf("query offers: %w", err)
	}
	defer rows.Close()

	var offers []entity.Offer
	for rows.Next() {
		var o entity.Offer
		if err := rows.Scan(&o.SellerID, &o.SellerName, &o.Description, &o.Image, &o.Price, &o.Quantity); err != nil {
			return make([]entity.Offer, 0), err
		}
		offers = append(offers, o)
	}

	return offers, nil
}

func (r *SellerRepo) CreateOffer(ctx context.Context, params *entity.CreateOfferParams) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("cannot start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var rootCategoryID int64
	err = tx.QueryRow(ctx, `
        SELECT categories_id FROM categories 
        WHERE name = $1 AND parent_id IS NULL`, params.CategoryName).Scan(&rootCategoryID)

	if errors.Is(err, sql.ErrNoRows) {
		err = tx.QueryRow(ctx, `
            INSERT INTO categories (name, parent_id) 
            VALUES ($1, NULL) 
            RETURNING categories_id`, params.CategoryName).Scan(&rootCategoryID)
		log.Println(err)
		if err != nil {
			return fmt.Errorf("create root category: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("find root category: %w", err)
	}

	var subCategoryID int64
	err = tx.QueryRow(ctx, `
        SELECT categories_id FROM categories 
        WHERE name = $1 AND parent_id = $2`, params.SubCategoryName, rootCategoryID).Scan(&subCategoryID)

	if errors.Is(err, sql.ErrNoRows) {
		err = tx.QueryRow(ctx, `
            INSERT INTO categories (name, parent_id) 
            VALUES ($1, $2) 
            RETURNING categories_id`, params.SubCategoryName, rootCategoryID).Scan(&subCategoryID)
		if err != nil {
			return fmt.Errorf("create subcategory: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("find subcategory: %w", err)
	}

	var productID int64
	createProdStmt := `
        INSERT INTO products (name, description, img, categories_id, created_at) 
        VALUES ($1, $2, $3, $4, NOW()) 
        RETURNING product_id`

	err = tx.QueryRow(ctx, createProdStmt,
		params.ProductName,
		params.Description,
		params.Image,
		subCategoryID,
	).Scan(&productID)
	if err != nil {
		return fmt.Errorf("create product: %w", err)
	}

	_, err = tx.Exec(ctx, `
        INSERT INTO seller_offers (seller_id, product_id, price, quantity) 
        VALUES ($1, $2, $3, $4)`,
		params.SellerID, productID, params.Price, params.Quantity)
	if err != nil {
		return fmt.Errorf("create offer: %w", err)
	}

	return tx.Commit(ctx)
}

func (r *SellerRepo) CreateOfferByExistProducts(ctx context.Context, params *entity.OfferWithID) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("cannot start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var exists bool

	err = tx.QueryRow(ctx, `
    SELECT EXISTS(
        SELECT 1 FROM categories 
        WHERE categories_id = $1 AND parent_id IS NULL
    )
`, params.CategoryID).Scan(&exists)

	if err != nil {
		log.Println(err)
		return fmt.Errorf("error checking existence of category: %w", err)
	}

	if !exists {
		return domain.ErrCategoryID
	}

	err = tx.QueryRow(ctx, `
    SELECT EXISTS(
        SELECT 1 FROM categories 
        WHERE categories_id = $1 AND parent_id = $2 )`,
		params.SubCategoryID, params.CategoryID,
	).Scan(&exists)

	if err != nil {
		log.Println(err)
		return fmt.Errorf("error checking existence of subcategory: %w", err)
	}

	if !exists {
		return domain.ErrSubCategoryID
	}

	err = tx.QueryRow(ctx, `
    SELECT EXISTS(
        SELECT 1 FROM products 
        WHERE categories_id = $1 AND product_id = $2 )`,
		params.SubCategoryID, params.ProductID,
	).Scan(&exists)

	if err != nil {
		log.Println(err)
		return fmt.Errorf("error checking existence of product: %w", err)
	}

	if !exists {
		return domain.ErrProductID
	}

	stmt := `INSERT INTO seller_offers (seller_id, product_id, price, quantity)
VALUES ($1, $2, $3, $4)
ON CONFLICT (seller_id, product_id)
DO UPDATE SET
    price = EXCLUDED.price,
    quantity = EXCLUDED.quantity;`

	_, err = tx.Exec(ctx, stmt, params.SellerID, params.ProductID, params.Price, params.Quantity)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("create offer: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("create offer: %w", err)
	}
	return nil
}
