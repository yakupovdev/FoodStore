package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yakupovdev/FoodStore/internal/delivery/http/dto"
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

func (h *SellerHandler) CreateOffer(ctx *gin.Context) {
	uid, ok := extractUserID(ctx)
	if !ok {
		ctx.JSON(ErrInvalidToken.Status, ErrInvalidToken.Response)
		return
	}

	var req dto.CreateOfferInput
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(ErrInvalidJSON.Status, ErrInvalidJSON.Response)
		return
	}

	output, err := h.uc.CreateOffer(ctx.Request.Context(), req, uid)
	if err != nil {
		ctx.JSON(ErrInternal.Status, ErrInternal.Response)
		return
	}

	ctx.JSON(http.StatusCreated, output)
}
