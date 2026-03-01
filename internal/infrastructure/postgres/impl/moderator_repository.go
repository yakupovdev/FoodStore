package impl

import (
	"context"
	"errors"
	"log"

	pg "github.com/jackc/pgx/v5"
	"github.com/yakupovdev/FoodStore/internal/domain"
	"github.com/yakupovdev/FoodStore/internal/domain/entity"
)

type ModeratorRepo struct {
	conn *pg.Conn
}

func NewModeratorRepo(conn *pg.Conn) *ModeratorRepo {
	return &ModeratorRepo{conn: conn}
}

func (r *ModeratorRepo) GetModerationSellerOffers(ctx context.Context) ([]entity.ModerationOffer, error) {
	stmt := `SELECT * FROM moderation_seller_offers`

	rows, err := r.conn.Query(ctx, stmt)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	var moderationOffers []entity.ModerationOffer

	for rows.Next() {
		var offer entity.ModerationOffer
		err = rows.Scan(&offer.SellerID,
			&offer.SellerName,
			&offer.SellerEmail,
			&offer.CategoryID,
			&offer.CategoryName,
			&offer.SubCategoryID,
			&offer.SubCategoryName,
			&offer.ProductID,
			&offer.ProductName,
			&offer.Description,
			&offer.Image,
			&offer.Price,
			&offer.Quantity)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		moderationOffers = append(moderationOffers, offer)
	}

	return moderationOffers, nil
}

func (r *ModeratorRepo) GetSellerOfferByProductID(ctx context.Context, productID int64) (*entity.ModerationOffer, error) {
	stmt := `SELECT * FROM moderation_seller_offers WHERE product_id = $1`

	var offer entity.ModerationOffer

	err := r.conn.QueryRow(ctx, stmt, productID).Scan(&offer.SellerID,
		&offer.SellerName,
		&offer.SellerEmail,
		&offer.CategoryID,
		&offer.CategoryName,
		&offer.SubCategoryID,
		&offer.SubCategoryName,
		&offer.ProductID,
		&offer.ProductName,
		&offer.Description,
		&offer.Image,
		&offer.Price,
		&offer.Quantity,
	)

	if err != nil {
		log.Println(err)
		if errors.Is(err, pg.ErrNoRows) {
			return nil, domain.ErrOfferNotFound
		}
		return nil, err
	}

	return &offer, nil
}

func (r *ModeratorRepo) DeleteModerationSellerOffer(ctx context.Context, productID int64) error {
	stmt := `DELETE FROM moderation_seller_offers WHERE product_id=$1`

	_, err := r.conn.Exec(ctx, stmt, productID)
	if err != nil {
		log.Println(err)
		if errors.Is(err, pg.ErrNoRows) {
			return domain.ErrOfferNotFound
		}
		return err
	}
	return nil
}

func (r *ModeratorRepo) CreateModerationOffer(ctx context.Context, params *entity.ModerationOffer) error {
	stmt := `INSERT INTO moderation_seller_offers (
                               seller_id, 
                               seller_name,
                               seller_email,
                               category_id,  
                               category_name, 
                               subcategory_id,
                               subcategory_name,
                               product_name, 
                               description, 
                               image, 
                               price, 
                               quantity) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	_, err := r.conn.Exec(ctx, stmt,
		params.SellerID,
		params.SellerName,
		params.SellerEmail,
		params.CategoryID,
		params.CategoryName,
		params.SubCategoryID,
		params.SubCategoryName,
		params.ProductName,
		params.Description,
		params.Image,
		params.Price,
		params.Quantity,
	)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
