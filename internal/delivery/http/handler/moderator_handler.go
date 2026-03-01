package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yakupovdev/FoodStore/internal/delivery/http/dto"
	"github.com/yakupovdev/FoodStore/internal/domain"
	"github.com/yakupovdev/FoodStore/internal/usecase"
)

type ModeratorHandler struct {
	uc *usecase.ModeratorUsecase
}

func NewModeratorHandler(uc *usecase.ModeratorUsecase) *ModeratorHandler {
	return &ModeratorHandler{uc: uc}
}

func (h *ModeratorHandler) GetExistProducts(ctx *gin.Context) {
	products, err := h.uc.GetAllExistProducts(ctx.Request.Context())
	if err != nil {
		ctx.JSON(ErrInternal.Status, ErrInternal.Response)
		return
	}

	ctx.JSON(http.StatusOK, products)
}

func (h *ModeratorHandler) GetModerationSellerOffers(ctx *gin.Context) {
	offers, err := h.uc.GetModerationSellerOffers(ctx.Request.Context())
	if err != nil {
		ctx.JSON(ErrInternal.Status, ErrInternal.Response)
		return
	}

	ctx.JSON(http.StatusOK, offers)
}

func (h *ModeratorHandler) ApproveOffer(ctx *gin.Context) {
	var req dto.OfferModerationAnswerInput

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(ErrInvalidJSON.Status, ErrInvalidJSON.Response)
		return
	}

	if req.ProductID <= 0 {
		ctx.JSON(ErrInvalidData.Status, ErrInvalidData.Response)
		return
	}

	res, err := h.uc.ApproveOffer(ctx, req)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrOfferNotFound):
			ctx.JSON(ErrOfferNotFound.Status, ErrOfferNotFound.Response)
		default:
			ctx.JSON(ErrInternal.Status, ErrInternal.Response)
		}
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (h *ModeratorHandler) RejectOffer(ctx *gin.Context) {
	var req dto.OfferModerationAnswerInput

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(ErrInvalidJSON.Status, ErrInvalidJSON.Response)
		return
	}

	if req.ProductID <= 0 {
		ctx.JSON(ErrInvalidData.Status, ErrInvalidData.Response)
		return
	}

	res, err := h.uc.RejectOffer(ctx, req)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrOfferNotFound):
			ctx.JSON(ErrOfferNotFound.Status, ErrOfferNotFound.Response)
		default:
			ctx.JSON(ErrInternal.Status, ErrInternal.Response)
		}
		return
	}

	ctx.JSON(http.StatusOK, res)
}
