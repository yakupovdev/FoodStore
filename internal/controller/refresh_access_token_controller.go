package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yakupovdev/FoodStore/internal/model"
	"github.com/yakupovdev/FoodStore/internal/usecase"
)

type RefreshAccessTokenController struct {
	ratu *usecase.RefreshAccessTokenUsecase
}

func NewRefreshAccessTokenController(ratu *usecase.RefreshAccessTokenUsecase) *RefreshAccessTokenController {
	return &RefreshAccessTokenController{ratu: ratu}
}

func (ratc *RefreshAccessTokenController) RefreshAccessToken(ctx *gin.Context) {
	userId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(ErrUserNotFound.Status, ErrUserNotFound.Response)
		return
	}
	userType, exists := ctx.Get("user_type")
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

	var utype string
	switch v := userType.(type) {
	case string:
		utype = v
	default:
		ctx.JSON(ErrInternal.Status, ErrInternal.Response)
	}

	accessToken, err := ratc.ratu.RefreshAccessToken(uid, utype)
	if err != nil {
		ctx.JSON(ErrInternal.Status, ErrInternal.Response)
		return
	}

	ctx.JSON(http.StatusOK, model.RefreshAccessTokenResponse{AccessToken: accessToken})
}
