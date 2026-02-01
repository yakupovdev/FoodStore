package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yakupovdev/FoodStore/internal/model"
	"github.com/yakupovdev/FoodStore/internal/usecase"
)

type SellerController struct {
	uc *usecase.SellerUsecase
}

func NewSellerController(uc *usecase.SellerUsecase) *SellerController {
	return &SellerController{
		uc: uc,
	}
}

func (c *SellerController) GetProfile(ctx *gin.Context) {
	userID, exist := ctx.Get("user_id")
	if !exist {
		ctx.JSON(ErrInvalidToken.Status, ErrInvalidToken.Response)
		return
	}

	var uid int64
	switch v := userID.(type) {
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

	resp, err := c.uc.GetProfileByID(uid)
	if err != nil {
		ctx.JSON(ErrInternal.Status, ErrInternal.Response)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (c *SellerController) GetOffers(ctx *gin.Context) {
	userID, exist := ctx.Get("user_id")
	if !exist {
		ctx.JSON(ErrInvalidToken.Status, ErrInvalidToken.Response)
		return
	}

	var uid int64
	switch v := userID.(type) {
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

	resp, err := c.uc.GetSellerOffersByID(uid)
	if err != nil {
		ctx.JSON(ErrInternal.Status, ErrInternal.Response)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (c *SellerController) CreateOffer(ctx *gin.Context) {
	userID, exist := ctx.Get("user_id")
	if !exist {
		ctx.JSON(ErrInvalidToken.Status, ErrInvalidToken.Response)
		return
	}

	var uid int64
	switch v := userID.(type) {
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

	var req model.CreateSellerOfferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(ErrInvalidJSON.Status, ErrInvalidJSON.Response)
		return
	}

	err := c.uc.CreateSellerOffer(req, uid)
	if err != nil {
		ctx.JSON(ErrInternal.Status, ErrInternal.Response)
		return
	}

	ctx.JSON(http.StatusCreated, model.CreateSellerOfferResponse{
		Message:         "Created successfully",
		ProductName:     req.ProductName,
		Description:     req.Description,
		Image:           req.Image,
		Price:           req.Price,
		Quantity:        req.Quantity,
		CategoryName:    req.CategoryName,
		SubCategoryName: req.SubCategoryName,
	})
}
