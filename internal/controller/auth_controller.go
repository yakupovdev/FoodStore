package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yakupovdev/FoodStore/internal/model"
	"github.com/yakupovdev/FoodStore/internal/usecase"
)

type AuthController struct {
	uc *usecase.AuthUsecase
}

func NewAuthController(uc *usecase.AuthUsecase) *AuthController {
	return &AuthController{uc: uc}
}

func (ac *AuthController) RegisterUser(ctx *gin.Context) {
	var req model.RegisterRequest

	err := ctx.ShouldBindJSON(&req)

	if err != nil {
		ctx.JSON(ErrInvalidJSON.Status, ErrInvalidJSON.Response)
		return
	}

	if req.Email == "" || req.Password == "" {
		ctx.JSON(ErrEmptyFields.Status, ErrEmptyFields.Response)
		return
	}

	err = ac.uc.RegisterUser(req.Email, req.Password, req.UserType, req.Balance)

	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrDuplicateEmail):
			ctx.JSON(ErrUserExists.Status, ErrUserExists.Response)
		default:
			ctx.JSON(ErrInternal.Status, ErrInternal.Response)
		}
		return
	}

	ctx.JSON(http.StatusOK, model.RegisterResponse{Message: "User Registered"})
}

func (ac *AuthController) LoginUser(ctx *gin.Context) {
	var req model.LoginRequest

	err := ctx.ShouldBindJSON(&req)

	if err != nil {
		ctx.JSON(ErrInvalidJSON.Status, ErrInvalidJSON.Response)
		return
	}

	if req.Email == "" || req.Password == "" {
		ctx.JSON(ErrEmptyFields.Status, ErrEmptyFields.Response)
		return
	}

	accessToken, refreshToken, err := ac.uc.LoginUser(req.Email, req.Password)

	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrInvalidCredentials):
			ctx.JSON(ErrInvalidCredentials.Status, ErrInvalidCredentials.Response)
		default:
			ctx.JSON(ErrInternal.Status, ErrInternal.Response)
		}
		return
	}

	ctx.JSON(http.StatusOK, model.LoginResponse{AccessToken: accessToken, RefreshToken: refreshToken})
}
