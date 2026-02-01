package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/yakupovdev/FoodStore/internal/model"
	"github.com/yakupovdev/FoodStore/internal/usecase"
)

type ClientController struct {
	uc *usecase.ClientUsecase
}

func NewClientController(uc *usecase.ClientUsecase) *ClientController {
	return &ClientController{uc: uc}
}

func (c *ClientController) GetProfile(ctx *gin.Context) {
	userID := ctx.GetInt64("user_id")
	userType := ctx.GetString("user_type")

	profile, err := c.uc.GetProfileByID(userID)
	if err != nil {
		ctx.JSON(ErrInternal.Status, ErrInternal.Response)
		return
	}
	var clientDTO model.ClientProfileDTO
	clientDTO.ID = profile.ID
	clientDTO.Name = profile.Name
	clientDTO.Email = profile.Email
	clientDTO.UserType = userType
	clientDTO.Balance = profile.Balance

	ctx.JSON(200, clientDTO)
}

func (c *ClientController) GetOrders(ctx *gin.Context) {
	userID := ctx.GetInt64("user_id")

	orders, err := c.uc.GetOrdersByClientID(userID)
	if err != nil {
		ctx.JSON(ErrInternal.Status, ErrInternal.Response)
		return
	}

	ctx.JSON(200, model.ClientOrdersDTO{
		OrderID:   orders[0].OrderID,
		ClientID:  orders[0].ClientID,
		Status:    orders[0].Status,
		CreatedAt: orders[0].CreatedAt,
		Items:     orders[0].Items,
	})
}
