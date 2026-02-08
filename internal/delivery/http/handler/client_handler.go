package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yakupovdev/FoodStore/internal/usecase"
)

type ClientHandler struct {
	uc *usecase.ClientUsecase
}

func NewClientHandler(uc *usecase.ClientUsecase) *ClientHandler {
	return &ClientHandler{uc: uc}
}

func (h *ClientHandler) GetProfile(ctx *gin.Context) {
	userID := ctx.GetInt64("user_id")
	userType := ctx.GetString("user_type")

	output, err := h.uc.GetProfileByID(ctx.Request.Context(), userID, userType)
	if err != nil {
		ctx.JSON(ErrInternal.Status, ErrInternal.Response)
		return
	}

	ctx.JSON(http.StatusOK, output)
}

func (h *ClientHandler) GetOrders(ctx *gin.Context) {
	userID := ctx.GetInt64("user_id")

	orders, err := h.uc.GetOrdersByClientID(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(ErrInternal.Status, ErrInternal.Response)
		return
	}

	if len(orders) == 0 {
		ctx.JSON(http.StatusOK, []interface{}{})
		return
	}

	ctx.JSON(http.StatusOK, orders[0])
}

func (h *ClientHandler) GetProducts(ctx *gin.Context) {
	categories, err := h.uc.GetProducts(ctx.Request.Context())
	if err != nil {
		ctx.JSON(ErrInternal.Status, ErrInternal.Response)
		return
	}

	ctx.JSON(http.StatusOK, categories)
}
