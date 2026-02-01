package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	pg "github.com/jackc/pgx/v5"
	"github.com/yakupovdev/FoodStore/internal/model"
)

type SellerRepository struct {
	Conn *pg.Conn
}

func NewSellerRepository(conn *pg.Conn) *SellerRepository {
	return &SellerRepository{
		Conn: conn,
	}
}

func (r *SellerRepository) GetSellerProfile(userID int64) (model.Seller, error) {
	ctx := context.Background()
	var seller model.Seller
	stmtBase := `
        SELECT u.email, u.balance, s.name, s.rating 
        FROM users u
        JOIN sellers s ON u.userid = s.seller_id
        WHERE u.userid = $1`

	err := r.Conn.QueryRow(ctx, stmtBase, userID).Scan(
		&seller.Email, &seller.Balance, &seller.Name, &seller.Rating,
	)
	if err != nil {
		return model.Seller{}, fmt.Errorf("find seller: %w", err)
	}

	seller.Type = "seller"

	return seller, nil
}

func (r *SellerRepository) GetSellerOffers(userID int64) ([]model.Offer, error) {
	ctx := context.Background()

	stmtOffers := `
        SELECT p.name, p.description, p.img, so.price, so.quantity
        FROM seller_offers so
        JOIN products p ON so.product_id = p.product_id
        WHERE so.seller_id = $1`

	rows, err := r.Conn.Query(ctx, stmtOffers, userID)
	if err != nil {
		return make([]model.Offer, 0), fmt.Errorf("query offers: %w", err)
	}
	defer rows.Close()

	var offers []model.Offer
	for rows.Next() {
		var o model.Offer
		if err := rows.Scan(&o.Name, &o.Description, &o.Image, &o.Price, &o.Quantity); err != nil {
			return make([]model.Offer, 0), err
		}
		offers = append(offers, o)
	}

	return offers, nil
}

func (r *SellerRepository) CreateSellerOffer(params model.CreateOfferParams) error {
	ctx := context.Background()

	tx, err := r.Conn.Begin(ctx)
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
