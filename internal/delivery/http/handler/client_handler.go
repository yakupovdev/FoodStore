package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yakupovdev/FoodStore/internal/delivery/http/dto"
	"github.com/yakupovdev/FoodStore/internal/domain"
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

func (h *ClientHandler) GetProductByID(ctx *gin.Context) {
	productIDStr := ctx.Param("product_id")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		log.Println(err, "invalid product ID in GetProductByID")
		ctx.JSON(ErrInvalidID.Status, ErrInvalidID.Response)
		return
	}

	output, err := h.uc.GetProductByID(ctx.Request.Context(), int64(productID))
	if err == domain.ErrNoProducts {
		ctx.JSON(ErrNoProducts.Status, ErrNoProducts.Response)
		return
	}
	if err != nil {
		log.Println(err, "error getting product by ID in GetProductByID")
		ctx.JSON(ErrInternal.Status, ErrInternal.Response)
		return
	}

	ctx.JSON(http.StatusOK, output)
}

func (h *ClientHandler) AddToCart(ctx *gin.Context) {
	var input dto.AddToCartInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		log.Println(err, "invalid JSON in AddToCart")
		ctx.JSON(ErrInvalidJSON.Status, ErrInvalidJSON.Response)
		return
	}
	clientID := ctx.GetInt64("user_id")
	input.ClientID = clientID

	output, err := h.uc.AddToCart(ctx.Request.Context(), input)
	if err != nil {
		log.Println(err, "error adding to cart in AddToCart")
		ctx.JSON(ErrInternal.Status, ErrInternal.Response)
		return
	}

	ctx.JSON(http.StatusOK, output)
}

func (h *ClientHandler) GetCartItems(ctx *gin.Context) {
	clientID := ctx.GetInt64("user_id")

	items, err := h.uc.GetCartItems(ctx.Request.Context(), clientID)
	if err != nil {
		log.Println(err, "error getting cart items in GetCartItems")
		ctx.JSON(ErrInternal.Status, ErrInternal.Response)
		return
	}

	ctx.JSON(http.StatusOK, items)
}

func (h *ClientHandler) CreateOrder(ctx *gin.Context) {
	var input dto.CreateOrderInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		log.Println(err, "invalid JSON in CreateOrder")
		ctx.JSON(ErrInvalidJSON.Status, ErrInvalidJSON.Response)
		return
	}
	clientID := ctx.GetInt64("user_id")
	input.ClientID = clientID

	output, err := h.uc.CreateOrder(ctx.Request.Context(), input)
	if err != nil {
		log.Println(err, "error creating order in CreateOrder")
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

	ctx.JSON(http.StatusOK, orders)
}

func (h *ClientHandler) GetProducts(ctx *gin.Context) {
	categories, err := h.uc.GetProducts(ctx.Request.Context())
	if err != nil {
		ctx.JSON(ErrInternal.Status, ErrInternal.Response)
		return
	}

	ctx.JSON(http.StatusOK, categories)
}

func (h *ClientHandler) UpdateBalance(ctx *gin.Context) {
	var input dto.UpdateBalanceInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(ErrInvalidJSON.Status, ErrInvalidJSON.Response)
		return
	}

	clientID := ctx.GetInt64("user_id")
	input.ClientID = clientID

	balanceUpdated, err := h.uc.UpdateBalance(ctx.Request.Context(), input)
	if err != nil {
		ctx.JSON(ErrInternal.Status, ErrInternal.Response)
		return
	}

	ctx.JSON(http.StatusOK, balanceUpdated)
}

func (h *ClientHandler) AddAdress(ctx *gin.Context) {
	var input dto.AddAddressInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(ErrInvalidJSON.Status, ErrInvalidJSON.Response)
		return
	}

	clientID := ctx.GetInt64("user_id")
	input.ClientID = clientID

	output, err := h.uc.AddAddress(ctx.Request.Context(), input)
	if err != nil {
		ctx.JSON(ErrInternal.Status, ErrInternal.Response)
		return
	}

	ctx.JSON(http.StatusOK, output)
}
