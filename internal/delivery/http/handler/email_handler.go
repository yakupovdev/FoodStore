package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yakupovdev/FoodStore/internal/domain"
	"github.com/yakupovdev/FoodStore/internal/usecase"
	"github.com/yakupovdev/FoodStore/internal/usecase/dto"
)

type EmailHandler struct {
	uc *usecase.RecoveryUsecase
}

func NewEmailHandler(uc *usecase.RecoveryUsecase) *EmailHandler {
	return &EmailHandler{uc: uc}
}

func (h *EmailHandler) SendCodeByEmail(ctx *gin.Context) {
	var req dto.SendCodeInput

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(ErrInvalidJSON.Status, ErrInvalidJSON.Response)
		return
	}

	if req.Email == "" || req.UserType == "" {
		ctx.JSON(ErrEmptyFields.Status, ErrEmptyFields.Response)
		return
	}

	err := h.uc.SendCodeByEmail(ctx.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			ctx.JSON(ErrUserNotFound.Status, ErrUserNotFound.Response)
		default:
			ctx.JSON(ErrInternal.Status, ErrInternal.Response)
		}
		return
	}

	ctx.JSON(http.StatusOK, dto.SendCodeOutput{Message: "Code sent successfully"})
}

func (h *EmailHandler) VerifyCode(ctx *gin.Context) {
	var req dto.VerifyCodeInput

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(ErrInvalidJSON.Status, ErrInvalidJSON.Response)
		return
	}

	if req.Email == "" || req.Code == "" || req.UserType == "" {
		ctx.JSON(ErrEmptyFields.Status, ErrEmptyFields.Response)
		return
	}

	output, err := h.uc.VerifyCode(ctx.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			ctx.JSON(ErrUserNotFound.Status, ErrUserNotFound.Response)
		case errors.Is(err, domain.ErrCodeInvalid):
			ctx.JSON(ErrVerifyCodeIsNotValid.Status, ErrVerifyCodeIsNotValid.Response)
		default:
			ctx.JSON(ErrInternal.Status, ErrInternal.Response)
		}
		return
	}

	ctx.JSON(http.StatusOK, output)
}
