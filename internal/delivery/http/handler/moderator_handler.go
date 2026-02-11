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

// GetExistProducts godoc
// @Summary Получить существующие продукты
// @Description Возвращает список всех существующих продуктов для модератора
// @Tags Moderator
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.CategoryIDOutput
// @Failure 500 {object} dto.ErrorOutput
// @Router /protected/moderator/products [get]
func (h *ModeratorHandler) GetExistProducts(ctx *gin.Context) {
	products, err := h.uc.GetAllExistProducts(ctx.Request.Context())
	if err != nil {
		ctx.JSON(ErrInternal.Status, ErrInternal.Response)
		return
	}

	ctx.JSON(http.StatusOK, products)
}

// GetModerationSellerOffers godoc
// @Summary Получить предложения на модерацию
// @Description Возвращает список предложений продавцов, ожидающих модерацию
// @Tags Moderator
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.OfferModerationOutput
// @Failure 500 {object} dto.ErrorOutput
// @Router /protected/moderator/offers [get]
func (h *ModeratorHandler) GetModerationSellerOffers(ctx *gin.Context) {
	offers, err := h.uc.GetModerationSellerOffers(ctx.Request.Context())
	if err != nil {
		ctx.JSON(ErrInternal.Status, ErrInternal.Response)
		return
	}

	ctx.JSON(http.StatusOK, offers)
}

// ApproveOffer godoc
// @Summary Одобрить предложение
// @Description Одобряет предложение продавца на модерации
// @Tags Moderator
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body dto.OfferModerationAnswerInput true "Данные одобрения"
// @Success 200 {object} dto.OfferModerationAnswerOutput
// @Failure 400 {object} dto.ErrorOutput "Некорректные данные"
// @Failure 404 {object} dto.ErrorOutput "Предложение не найдено"
// @Failure 500 {object} dto.ErrorOutput
// @Router /protected/moderator/approve [post]
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

// RejectOffer godoc
// @Summary Отклонить предложение
// @Description Отклоняет предложение продавца на модерации
// @Tags Moderator
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body dto.OfferModerationAnswerInput true "Данные отклонения"
// @Success 200 {object} dto.OfferModerationAnswerOutput
// @Failure 400 {object} dto.ErrorOutput "Некорректные данные"
// @Failure 404 {object} dto.ErrorOutput "Предложение не найдено"
// @Failure 500 {object} dto.ErrorOutput
// @Router /protected/moderator/reject [post]
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
