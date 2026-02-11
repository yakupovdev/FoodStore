package repository

import (
	"context"

	"github.com/yakupovdev/FoodStore/internal/domain/entity"
)

type SellerRepository interface {
	FindByUserID(ctx context.Context, userID int64) (*entity.Seller, error)

	GetOffersBySellerID(ctx context.Context, sellerID int64) ([]entity.Offer, error)

	GetOffersByProductID(ctx context.Context, productID int64) ([]entity.Offer, error)

	CreateOffer(ctx context.Context, params *entity.CreateOfferParams) error

	CreateOfferByExistProducts(ctx context.Context, params *entity.OfferWithID) error

	GetOffersBySellerIDAndProductID(ctx context.Context, sellerID, productID int64) (*entity.Offer, error)

	DeleteOffer(ctx context.Context, params *entity.OfferPrimary) error

	UpdateOffer(ctx context.Context, params *entity.SellerOffer) error

	DecreaseOfferQuantity(ctx context.Context, params *entity.OfferQuantity) error

	PurchaseSubscription(ctx context.Context, params *entity.PurchaseSubscriptionParams) error
}
