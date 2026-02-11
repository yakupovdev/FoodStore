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

// GetProfile godoc
// @Summary Получить профиль клиента
// @Description Возвращает профиль текущего авторизованного клиента
// @Tags Client
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.ClientProfileOutput
// @Failure 500 {object} dto.ErrorOutput
// @Router /protected/client/profile [get]
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

// GetProductByID godoc
// @Summary Получить продукт по ID
// @Description Возвращает информацию о продукте по его ID с предложениями продавцов
// @Tags Client
// @Produce json
// @Security BearerAuth
// @Param product_id path int true "ID продукта"
// @Success 200 {object} dto.ProductOutput
// @Failure 400 {object} dto.ErrorOutput "Некорректный ID"
// @Failure 404 {object} dto.ErrorOutput "Продукт не найден"
// @Failure 500 {object} dto.ErrorOutput
// @Router /protected/client/products/{product_id} [get]
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

// AddToCart godoc
// @Summary Добавить товар в корзину
// @Description Добавляет товар в корзину клиента
// @Tags Client
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body dto.AddToCartInput true "Данные для добавления в корзину"
// @Success 200 {object} dto.AddToCartOutput
// @Failure 400 {object} dto.ErrorOutput
// @Failure 500 {object} dto.ErrorOutput
// @Router /protected/client/cart [post]
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

// GetCartItems godoc
// @Summary Получить содержимое корзины
// @Description Возвращает все товары в корзине клиента
// @Tags Client
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.CartItemOutput
// @Failure 500 {object} dto.ErrorOutput
// @Router /protected/client/cart [get]
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

// CreateOrder godoc
// @Summary Создать заказ
// @Description Создаёт новый заказ клиента
// @Tags Client
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body dto.CreateOrderInput true "Данные заказа"
// @Success 200 {object} dto.CreateOrderOutput
// @Failure 400 {object} dto.ErrorOutput "Некорректный JSON или недостаточно товара"
// @Failure 404 {object} dto.ErrorOutput "Продукт не найден"
// @Failure 500 {object} dto.ErrorOutput
// @Router /protected/client/orders [post]
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
		if err == domain.ErrNoProducts {
			ctx.JSON(ErrNoProducts.Status, ErrNoProducts.Response)
			return
		}
		if err == domain.ErrNotEnoughQuantity {
			ctx.JSON(ErrNotEnoughQuantity.Status, ErrNotEnoughQuantity.Response)
			return
		}
		ctx.JSON(ErrInternal.Status, ErrInternal.Response)
		return
	}

	ctx.JSON(http.StatusOK, output)
}

// GetOrders godoc
// @Summary Получить заказы клиента
// @Description Возвращает список всех заказов текущего клиента
// @Tags Client
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.ClientOrderOutput
// @Failure 500 {object} dto.ErrorOutput
// @Router /protected/client/orders [get]
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

// GetProducts godoc
// @Summary Получить каталог продуктов
// @Description Возвращает полный каталог продуктов по категориям с предложениями продавцов
// @Tags Client
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.CategoryOutput
// @Failure 404 {object} dto.ErrorOutput "Продукты не найдены"
// @Failure 500 {object} dto.ErrorOutput
// @Router /protected/client/products [get]
func (h *ClientHandler) GetProducts(ctx *gin.Context) {
	categories, err := h.uc.GetProducts(ctx.Request.Context())
	if err != nil {
		if err == domain.ErrNoProducts {
			ctx.JSON(ErrNoProducts.Status, ErrNoProducts.Response)
			return
		}
		if err == domain.ErrOfferNotFound {
			ctx.JSON(ErrOfferNotFound.Status, ErrOfferNotFound.Response)
			return
		}
		log.Println(err)
		ctx.JSON(ErrInternal.Status, ErrInternal.Response)
		return
	}

	ctx.JSON(http.StatusOK, categories)
}

// AddAdress godoc
// @Summary Добавить адрес клиента
// @Description Добавляет новый адрес доставки клиенту
// @Tags Client
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body dto.AddAddressInput true "Данные адреса"
// @Success 200 {object} dto.AddAddressOutput
// @Failure 400 {object} dto.ErrorOutput
// @Failure 500 {object} dto.ErrorOutput
// @Router /protected/client/address [post]
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

// GetProductsByPriority godoc
// @Summary Получить продукты по приоритету
// @Description Возвращает продукты, отсортированные по приоритету (продавцы с подпиской выше)
// @Tags Client
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.CategoryOutput
// @Failure 404 {object} dto.ErrorOutput "Продукты не найдены"
// @Failure 500 {object} dto.ErrorOutput
// @Router /protected/client/ [get]
func (h *ClientHandler) GetProductsByPriority(ctx *gin.Context) {
	products, err := h.uc.GetProductsByPriority(ctx.Request.Context())
	if err != nil {
		if err == domain.ErrNoProducts {
			ctx.JSON(ErrNoProducts.Status, ErrNoProducts.Response)
			return
		}
		log.Println(err)
		ctx.JSON(ErrInternal.Status, ErrInternal.Response)
		return
	}

	ctx.JSON(http.StatusOK, products)
}
