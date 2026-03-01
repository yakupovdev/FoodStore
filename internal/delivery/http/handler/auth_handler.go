package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yakupovdev/FoodStore/internal/delivery/http/dto"
	"github.com/yakupovdev/FoodStore/internal/domain"
	"github.com/yakupovdev/FoodStore/internal/usecase"
)

type AuthHandler struct {
	uc *usecase.AuthUsecase
}

func NewAuthHandler(uc *usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{uc: uc}
}

func (h *AuthHandler) RegisterUser(ctx *gin.Context) {
	var req dto.RegisterInput

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(ErrInvalidJSON.Status, ErrInvalidJSON.Response)
		return
	}

	if req.Email == "" || req.Password == "" || req.UserType == "" || req.Name == "" {
		ctx.JSON(ErrEmptyFields.Status, ErrEmptyFields.Response)
		return
	}

	if lowerUserType := strings.ToLower(req.UserType); lowerUserType != "client" && lowerUserType != "seller" && lowerUserType != "moderator" && lowerUserType != "admin" {
		ctx.JSON(ErrInvalidUserType.Status, ErrInvalidUserType.Response)
		return
	}

	output, err := h.uc.Register(ctx.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserAlreadyExists):
			ctx.JSON(ErrUserExists.Status, ErrUserExists.Response)
		default:
			ctx.JSON(ErrInternal.Status, ErrInternal.Response)
		}
		return
	}

	ctx.JSON(http.StatusOK, output)
}

func (h *AuthHandler) LoginUser(ctx *gin.Context) {
	var req dto.LoginInput

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(ErrInvalidJSON.Status, ErrInvalidJSON.Response)
		return
	}

	if req.Email == "" || req.Password == "" || req.UserType == "" {
		ctx.JSON(ErrEmptyFields.Status, ErrEmptyFields.Response)
		return
	}

	output, err := h.uc.Login(ctx.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidCredentials):
			ctx.JSON(ErrInvalidCredentials.Status, ErrInvalidCredentials.Response)
		default:
			ctx.JSON(ErrInternal.Status, ErrInternal.Response)
		}
		return
	}

	ctx.JSON(http.StatusOK, output)
}
