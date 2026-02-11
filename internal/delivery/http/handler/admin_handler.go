package handler

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/yakupovdev/FoodStore/internal/delivery/http/dto"
	"github.com/yakupovdev/FoodStore/internal/usecase"
)

type AdminHandler struct {
	uc *usecase.AdminUsecase
}

func NewAdminHandler(uc *usecase.AdminUsecase) *AdminHandler {
	return &AdminHandler{uc: uc}
}

func (h *AdminHandler) DeleteUser(ctx *gin.Context) {
	var req dto.DeleteUserInput
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(ErrInvalidJSON.Status, ErrInvalidJSON.Response)
		return
	}

	output, err := h.uc.DeleteUser(ctx.Request.Context(), req)
	if err != nil {
		log.Println(err)
		ctx.JSON(ErrInternal.Status, ErrInternal.Response)
		return
	}
	ctx.JSON(204, output)
}

func (h *AdminHandler) UpdateBalance(ctx *gin.Context) {
	var req dto.UpdateBalanceInput
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(ErrInvalidJSON.Status, ErrInvalidJSON.Response)
		return
	}

	output, err := h.uc.UpdateBalance(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(ErrInternal.Status, ErrInternal.Response)
		return
	}
	ctx.JSON(200, output)
}

func (h *AdminHandler) GetAllLogTransactions(ctx *gin.Context) {
	output, err := h.uc.GetAllLogTransactions(ctx.Request.Context())
	if err != nil {
		ctx.JSON(ErrInternal.Status, ErrInternal.Response)
		return
	}
	ctx.JSON(200, output)
}

func (h *AdminHandler) GetAllUsers(ctx *gin.Context) {
	output, err := h.uc.GetAllUsers(ctx.Request.Context())
	if err != nil {
		ctx.JSON(ErrInternal.Status, ErrInternal.Response)
		return
	}
	ctx.JSON(200, output)
}

func (h *AdminHandler) CreateAdmin(ctx *gin.Context) {
	var req dto.CreateAdminInput
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(ErrInvalidJSON.Status, ErrInvalidJSON.Response)
		return
	}
	output, err := h.uc.CreateAdmin(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(ErrInternal.Status, ErrInternal.Response)
		return
	}
	ctx.JSON(201, output)
}
