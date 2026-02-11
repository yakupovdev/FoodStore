package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yakupovdev/FoodStore/internal/delivery/http/dto"
	"github.com/yakupovdev/FoodStore/internal/domain"
	"github.com/yakupovdev/FoodStore/internal/usecase"
)

type SellerHandler struct {
	uc *usecase.SellerUsecase
}

func NewSellerHandler(uc *usecase.SellerUsecase) *SellerHandler {
	return &SellerHandler{uc: uc}
}

func (h *SellerHandler) GetProfile(ctx *gin.Context) {
	uid, ok := extractUserID(ctx)
	if !ok {
		ctx.JSON(ErrInvalidToken.Status, ErrInvalidToken.Response)
		return
	}

	output, err := h.uc.GetProfileByID(ctx.Request.Context(), uid)
	if err != nil {
		ctx.JSON(ErrInternal.Status, ErrInternal.Response)
		return
	}

	ctx.JSON(http.StatusOK, output)
}

func (h *SellerHandler) GetOffers(ctx *gin.Context) {
	uid, ok := extractUserID(ctx)
	if !ok {
		ctx.JSON(ErrInvalidToken.Status, ErrInvalidToken.Response)
		return
	}

	output, err := h.uc.GetOffersByID(ctx.Request.Context(), uid)
	if err != nil {
		ctx.JSON(ErrInternal.Status, ErrInternal.Response)
		return
	}

	ctx.JSON(http.StatusOK, output)
}

func (h *SellerHandler) GetExistProducts(ctx *gin.Context) {
	products, err := h.uc.GetAllExistProducts(ctx.Request.Context())
	if err != nil {
		ctx.JSON(ErrInternal.Status, ErrInternal.Response)
		return
	}

	ctx.JSON(http.StatusOK, products)
}

func (h *SellerHandler) CreateOfferByExistProducts(ctx *gin.Context) {
	var req dto.CreateOfferByExistProductsInput

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(ErrInvalidJSON.Status, ErrInvalidJSON.Response)
		return
	}
	var ok bool
	req.SellerID, ok = extractUserID(ctx)
	if !ok {
		ctx.JSON(ErrInvalidToken.Status, ErrInvalidToken.Response)
		return
	}

	res, err := h.uc.CreateOfferByExistProducts(ctx, req)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrCategoryID) || errors.Is(err, domain.ErrSubCategoryID) || errors.Is(err, domain.ErrProductID) || errors.Is(err, domain.ErrInvalidPrice) || errors.Is(err, domain.ErrInvalidQuantity):
			ctx.JSON(ErrInvalidData.Status, ErrInvalidData.Response)
		default:
			ctx.JSON(ErrInternal.Status, ErrInternal.Response)
		}
	}

	ctx.JSON(http.StatusCreated, res)
}

func (h *SellerHandler) UpdateOffer(ctx *gin.Context) {
	uid, ok := extractUserID(ctx)
	if !ok {
		ctx.JSON(ErrInvalidToken.Status, ErrInvalidToken.Response)
		return
	}

	var req dto.UpdateOfferInput
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(ErrInvalidJSON.Status, ErrInvalidJSON.Response)
		return
	}

	req.SellerID = uid
	output, err := h.uc.UpdateOffer(ctx.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrProductID) || errors.Is(err, domain.ErrInvalidPrice) || errors.Is(err, domain.ErrInvalidQuantity):
			ctx.JSON(ErrInvalidData.Status, ErrInvalidData.Response)
		case errors.Is(err, domain.ErrOfferNotFound):
			ctx.JSON(ErrOfferNotFound.Status, ErrOfferNotFound.Response)
		default:
			ctx.JSON(ErrInternal.Status, ErrInternal.Response)
		}
		return
	}

	ctx.JSON(http.StatusOK, output)
}

func (h *SellerHandler) DeleteOffer(ctx *gin.Context) {
	uid, ok := extractUserID(ctx)
	if !ok {
		ctx.JSON(ErrInvalidToken.Status, ErrInvalidToken.Response)
		return
	}

	var req dto.DeleteOfferInput
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(ErrInvalidJSON.Status, ErrInvalidJSON.Response)
		return
	}

	req.SellerID = uid

	output, err := h.uc.DeleteOffer(ctx.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrProductID):
			ctx.JSON(ErrInvalidData.Status, ErrInvalidData.Response)
		case errors.Is(err, domain.ErrOfferNotFound):
			ctx.JSON(ErrOfferNotFound.Status, ErrOfferNotFound.Response)
		default:
			ctx.JSON(ErrInternal.Status, ErrInternal.Response)
		}
		return
	}

	ctx.JSON(http.StatusNoContent, output)
}

func (h *SellerHandler) CreateOfferWithNewProduct(ctx *gin.Context) {
	uid, ok := extractUserID(ctx)
	if !ok {
		ctx.JSON(ErrInvalidToken.Status, ErrInvalidToken.Response)
		return
	}

	var req dto.CreateOfferWithNewProductInput
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(ErrInvalidJSON.Status, ErrInvalidJSON.Response)
		return
	}

	req.SellerID = uid
	output, err := h.uc.CreateOfferWithNewProduct(ctx.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrSubCategoryID) ||
			errors.Is(err, domain.ErrSubCategoryNotFound) ||
			errors.Is(err, domain.ErrCategoryNotFound) ||
			errors.Is(err, domain.ErrInvalidPrice) ||
			errors.Is(err, domain.ErrInvalidQuantity) ||
			errors.Is(err, domain.ErrInvalidProductName) ||
			errors.Is(err, domain.ErrInvalidDescription) ||
			errors.Is(err, domain.ErrInvalidQuantity) ||
			errors.Is(err, domain.ErrInvalidPrice):
			ctx.JSON(ErrInvalidData.Status, ErrInvalidData.Response)
		default:
			ctx.JSON(ErrInternal.Status, ErrInternal.Response)
		}
		return
	}

	ctx.JSON(http.StatusCreated, output)
}
