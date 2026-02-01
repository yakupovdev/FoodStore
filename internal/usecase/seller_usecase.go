package usecase

import (
	"github.com/yakupovdev/FoodStore/internal/model"
	"github.com/yakupovdev/FoodStore/internal/repository"
)

type SellerUsecase struct {
	repo *repository.SellerRepository
}

func NewSellerUsecase(repo *repository.SellerRepository) (*SellerUsecase, error) {
	if repo == nil {
		return nil, ErrDatabaseConnection
	}
	return &SellerUsecase{
		repo: repo,
	}, nil
}

func (uc *SellerUsecase) GetProfileByID(userID int64) (model.SellerProfileResponse, error) {
	seller, err := uc.repo.GetSellerProfile(userID)
	if err != nil {
		return model.SellerProfileResponse{}, err
	}

	return model.SellerProfileResponse{
		Name:    seller.Name,
		Email:   seller.Email,
		Type:    seller.Type,
		Balance: seller.Balance,
		Rating:  seller.Rating,
	}, nil

}

func (uc *SellerUsecase) GetSellerOffersByID(userID int64) (model.SellerOffersResponse, error) {
	offers, err := uc.repo.GetSellerOffers(userID)
	if err != nil {
		return model.SellerOffersResponse{}, err
	}

	return model.SellerOffersResponse{
		Offers: offers,
	}, nil
}

func (uc *SellerUsecase) CreateSellerOffer(req model.CreateSellerOfferRequest, userID int64) error {
	err := uc.repo.CreateSellerOffer(model.CreateOfferParams{
		SellerID:        userID,
		ProductName:     req.ProductName,
		Description:     req.Description,
		Image:           req.Image,
		Price:           req.Price,
		Quantity:        req.Quantity,
		CategoryName:    req.CategoryName,
		SubCategoryName: req.SubCategoryName,
	})

	if err != nil {
		return err
	}
	return nil
}
