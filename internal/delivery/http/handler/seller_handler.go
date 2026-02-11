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

// GetProfile godoc
// @Summary Получить профиль продавца
// @Description Возвращает профиль текущего авторизованного продавца
// @Tags Seller
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.SellerProfileOutput
// @Failure 401 {object} dto.ErrorOutput "Невалидный токен"
// @Failure 500 {object} dto.ErrorOutput
// @Router /protected/seller/profile [get]
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

// GetOffers godoc
// @Summary Получить предложения продавца
// @Description Возвращает список всех предложений текущего продавца
// @Tags Seller
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.SellerOffersListOutput
// @Failure 401 {object} dto.ErrorOutput
// @Failure 500 {object} dto.ErrorOutput
// @Router /protected/seller/offers [get]
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

// GetExistProducts godoc
// @Summary Получить существующие продукты
// @Description Возвращает список всех существующих продуктов с категориями и подкатегориями
// @Tags Seller
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.CategoryIDOutput
// @Failure 500 {object} dto.ErrorOutput
// @Router /protected/seller/products [get]
func (h *SellerHandler) GetExistProducts(ctx *gin.Context) {
	products, err := h.uc.GetAllExistProducts(ctx.Request.Context())
	if err != nil {
		ctx.JSON(ErrInternal.Status, ErrInternal.Response)
		return
	}

	ctx.JSON(http.StatusOK, products)
}

// CreateOfferByExistProducts godoc
// @Summary Создать предложение на существующий продукт
// @Description Создаёт предложение продавца на уже существующий в каталоге продукт
// @Tags Seller
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body dto.CreateOfferByExistProductsInput true "Данные предложения"
// @Success 201 {object} dto.CreateOfferByExistProductsOutput
// @Failure 400 {object} dto.ErrorOutput "Некорректные данные"
// @Failure 401 {object} dto.ErrorOutput
// @Failure 500 {object} dto.ErrorOutput
// @Router /protected/seller/offers [post]
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

// UpdateOffer godoc
// @Summary Обновить предложение
// @Description Обновляет цену и количество в предложении продавца
// @Tags Seller
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body dto.UpdateOfferInput true "Данные для обновления"
// @Success 200 {object} dto.UpdateOfferOutput
// @Failure 400 {object} dto.ErrorOutput "Некорректные данные"
// @Failure 401 {object} dto.ErrorOutput
// @Failure 404 {object} dto.ErrorOutput "Предложение не найдено"
// @Failure 500 {object} dto.ErrorOutput
// @Router /protected/seller/offers [put]
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

// DeleteOffer godoc
// @Summary Удалить предложение
// @Description Удаляет предложение продавца по product_id
// @Tags Seller
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body dto.DeleteOfferInput true "Данные для удаления"
// @Success 204 {object} dto.DeleteOfferOutput
// @Failure 400 {object} dto.ErrorOutput "Некорректные данные"
// @Failure 401 {object} dto.ErrorOutput
// @Failure 404 {object} dto.ErrorOutput "Предложение не найдено"
// @Failure 500 {object} dto.ErrorOutput
// @Router /protected/seller/offers [delete]
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

// CreateOfferWithNewProduct godoc
// @Summary Создать предложение с новым продуктом
// @Description Создаёт новый продукт и предложение продавца на него (отправляется на модерацию)
// @Tags Seller
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body dto.CreateOfferWithNewProductInput true "Данные нового продукта и предложения"
// @Success 201 {object} dto.CreateOfferWithNewProductOutput
// @Failure 400 {object} dto.ErrorOutput "Некорректные данные"
// @Failure 401 {object} dto.ErrorOutput
// @Failure 500 {object} dto.ErrorOutput
// @Router /protected/seller/new-offers [post]
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

// PurchaseSubscription godoc
// @Summary Купить подписку
// @Description Покупка подписки продавца для приоритетного отображения товаров
// @Tags Seller
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.PurchaseSubscriptionOutput
// @Failure 400 {object} dto.ErrorOutput "Недостаточно средств или ошибка подписки"
// @Failure 401 {object} dto.ErrorOutput
// @Failure 500 {object} dto.ErrorOutput
// @Router /protected/seller/subscription [get]
func (h *SellerHandler) PurchaseSubscription(ctx *gin.Context) {
	uid, ok := extractUserID(ctx)
	if !ok {
		ctx.JSON(ErrInvalidToken.Status, ErrInvalidToken.Response)
		return
	}

	var req dto.PurchaseSubscriptionInput
	req = dto.PurchaseSubscriptionInput{
		ID: uid,
	}

	output, err := h.uc.PurchaseSubscription(ctx.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrSubscriptionNotFound):
			ctx.JSON(ErrSubscriptionNotFound.Status, ErrSubscriptionNotFound.Response)
		case errors.Is(err, domain.ErrNotEnoughBalance):
			ctx.JSON(ErrNotEnoughBalance.Status, ErrNotEnoughBalance.Response)
		default:
			ctx.JSON(ErrInternal.Status, ErrInternal.Response)
		}
		return
	}

	ctx.JSON(http.StatusOK, output)
}
