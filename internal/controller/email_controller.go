package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/yakupovdev/FoodStore/internal/model"
	"github.com/yakupovdev/FoodStore/usecase"
)

type EmailController struct {
	eu *usecase.EmailUsecase
}

func NewEmailController(eu *usecase.EmailUsecase) *EmailController {
	return &EmailController{eu: eu}
}

func (ec *EmailController) SendCodeByEmail(ctx *gin.Context) {
	var req model.EmailRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(ErrInvalidJSON.Status, ErrInvalidJSON.Response)
		return
	}

	if req.Email == "" {
		ctx.JSON(ErrEmptyFields.Status, ErrEmptyFields.Response)
		return
	}

	err = ec.eu.SendCodeByEmail(req.Email)
	if err != nil { // TODO: handle specific errors
		switch {
		default:
			ctx.JSON(ErrInternal.Status, ErrInternal.Response)
		}
	}

	ctx.JSON(200, model.EmailResponse{Message: "Code sent successfully"})

}
