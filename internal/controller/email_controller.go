package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yakupovdev/FoodStore/internal/model"
	"github.com/yakupovdev/FoodStore/internal/repository"
	"github.com/yakupovdev/FoodStore/usecase"
)

type EmailController struct {
	eu *usecase.EmailUsecase
}

func NewEmailController(eu *usecase.EmailUsecase) *EmailController {
	return &EmailController{eu: eu}
}

func (ec *EmailController) SendCodeByEmail(ctx *gin.Context) {
	var req model.VerifyEmailRequest

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
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrUserNotFound):
			ctx.JSON(ErrUserNotFound.Status, ErrUserNotFound.Response)
		default:

			ctx.JSON(ErrInternal.Status, ErrInternal.Response)
		}
		return
	}

	ctx.JSON(http.StatusOK, model.VerifyEmailResponse{Message: "Code sent successfully"})
}

func (ec *EmailController) VerifyCode(ctx *gin.Context) {
	var req model.VerifyCodeRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(ErrInvalidJSON.Status, ErrInvalidJSON.Response)
		return
	}

	if req.Email == "" || req.Code == "" {
		ctx.JSON(ErrEmptyFields.Status, ErrEmptyFields.Response)
		return
	}

	token, err := ec.eu.VerifyCode(req.Email, req.Code)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrUserNotFound):
			ctx.JSON(ErrUserNotFound.Status, ErrUserNotFound.Response)
		case errors.Is(err, usecase.ErrCodeIsNotValid):
			ctx.JSON(ErrVerifyCodeIsNotValid.Status, ErrVerifyCodeIsNotValid.Response)
		default:
			ctx.JSON(ErrInternal.Status, ErrInternal.Response)
		}
		return
	}

	ctx.JSON(http.StatusOK, model.VerifyCodeResponse{Token: token, Message: "Code verified successfully"})
}
