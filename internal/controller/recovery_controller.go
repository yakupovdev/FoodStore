package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yakupovdev/FoodStore/internal/model"
	"github.com/yakupovdev/FoodStore/usecase"
)

type RecoveryController struct {
	ru *usecase.RecoveryUsecase
}

func NewRecoveryController(ru *usecase.RecoveryUsecase) *RecoveryController {
	return &RecoveryController{ru}
}

func (rc *RecoveryController) ResetUserPassword(ctx *gin.Context) {
	var req model.ResetUserPasswordRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(ErrInvalidJSON.Status, ErrInvalidJSON.Response)
		return
	}

	if req.NewPassword == "" {
		ctx.JSON(ErrEmptyFields.Status, ErrEmptyFields.Response)
	}

	userId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(ErrUserNotFound.Status, ErrUserNotFound.Response)
		return
	}

	var uid int64
	switch v := userId.(type) {
	case int64:
		uid = v
	case int:
		uid = int64(v)
	case float64:
		uid = int64(v)
	case string:
		parsed, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			ctx.JSON(ErrInvalidJSON.Status, ErrInvalidJSON.Response)
			return
		}
		uid = parsed
	default:
		ctx.JSON(ErrInternal.Status, ErrInternal.Response)
		return
	}

	err = rc.ru.ResetUserPassword(uid, req.NewPassword)
	if err != nil {
		ctx.JSON(ErrInternal.Status, ErrInternal.Response)
		return
	}

	ctx.JSON(http.StatusOK, model.ResetUserPasswordResponse{Message: "Password reset"})
}
