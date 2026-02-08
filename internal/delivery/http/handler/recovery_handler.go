package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yakupovdev/FoodStore/internal/delivery/http/dto"
	"github.com/yakupovdev/FoodStore/internal/usecase"
)

type RecoveryHandler struct {
	uc *usecase.RecoveryUsecase
}

func NewRecoveryHandler(uc *usecase.RecoveryUsecase) *RecoveryHandler {
	return &RecoveryHandler{uc: uc}
}

func (h *RecoveryHandler) ResetUserPassword(ctx *gin.Context) {
	var req dto.ResetPasswordInput

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(ErrInvalidJSON.Status, ErrInvalidJSON.Response)
		return
	}

	if req.NewPassword == "" {
		ctx.JSON(ErrEmptyFields.Status, ErrEmptyFields.Response)
		return
	}

	uid, ok := extractUserID(ctx)
	if !ok {
		ctx.JSON(ErrUserNotFound.Status, ErrUserNotFound.Response)
		return
	}
	req.UserID = uid

	err := h.uc.ResetPassword(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(ErrInternal.Status, ErrInternal.Response)
		return
	}

	ctx.JSON(http.StatusOK, dto.ResetPasswordOutput{Message: "Password reset"})
}
