package repository

import (
	"context"

	"github.com/yakupovdev/FoodStore/internal/domain/entity"
)

type ModeratorRepository interface {
	GetModerationSellerOffers(ctx context.Context) ([]entity.ModerationOffer, error)

	DeleteModerationSellerOffer(ctx context.Context, productID int64) error

	GetSellerOfferByProductID(ctx context.Context, productID int64) (*entity.ModerationOffer, error)

	CreateModerationOffer(ctx context.Context, params *entity.ModerationOffer) error
}
