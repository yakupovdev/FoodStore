package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yakupovdev/FoodStore/internal/usecase"
)

type RefreshTokenHandler struct {
	uc *usecase.AuthUsecase
}

func NewRefreshTokenHandler(uc *usecase.AuthUsecase) *RefreshTokenHandler {
	return &RefreshTokenHandler{uc: uc}
}

func (h *RefreshTokenHandler) RefreshAccessToken(ctx *gin.Context) {
	uid, ok := extractUserID(ctx)
	if !ok {
		ctx.JSON(ErrUserNotFound.Status, ErrUserNotFound.Response)
		return
	}

	utype, ok := extractUserType(ctx)
	if !ok {
		ctx.JSON(ErrUserNotFound.Status, ErrUserNotFound.Response)
		return
	}

	output, err := h.uc.RefreshAccessToken(ctx.Request.Context(), uid, utype)
	if err != nil {
		ctx.JSON(ErrInternal.Status, ErrInternal.Response)
		return
	}

	ctx.JSON(http.StatusOK, output)
}
