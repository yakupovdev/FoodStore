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

// DeleteUser godoc
// @Summary Удалить пользователя
// @Description Удаляет пользователя по user_id
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body dto.DeleteUserInput true "Данные для удаления"
// @Success 204 {object} dto.DeleteUserOutput
// @Failure 400 {object} dto.ErrorOutput
// @Failure 500 {object} dto.ErrorOutput
// @Router /protected/admin/users [delete]
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

// UpdateBalance godoc
// @Summary Обновить баланс пользователя
// @Description Обновляет баланс пользователя по user_id
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body dto.UpdateBalanceInput true "Данные баланса"
// @Success 200 {object} dto.BalanceUpdateOutput
// @Failure 400 {object} dto.ErrorOutput
// @Failure 500 {object} dto.ErrorOutput
// @Router /protected/admin/balance [post]
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

// GetAllLogTransactions godoc
// @Summary Получить логи транзакций
// @Description Возвращает все логи финансовых транзакций
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.GetLogsHistoryOutput
// @Failure 500 {object} dto.ErrorOutput
// @Router /protected/admin/logs [get]
func (h *AdminHandler) GetAllLogTransactions(ctx *gin.Context) {
	output, err := h.uc.GetAllLogTransactions(ctx.Request.Context())
	if err != nil {
		ctx.JSON(ErrInternal.Status, ErrInternal.Response)
		return
	}
	ctx.JSON(200, output)
}

// GetAllUsers godoc
// @Summary Получить всех пользователей
// @Description Возвращает список всех пользователей системы
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.GetAllUsersOutput
// @Failure 500 {object} dto.ErrorOutput
// @Router /protected/admin/users [get]
func (h *AdminHandler) GetAllUsers(ctx *gin.Context) {
	output, err := h.uc.GetAllUsers(ctx.Request.Context())
	if err != nil {
		ctx.JSON(ErrInternal.Status, ErrInternal.Response)
		return
	}
	ctx.JSON(200, output)
}

// CreateAdmin godoc
// @Summary Создать администратора
// @Description Создаёт нового администратора по секретному ключу
// @Tags Admin
// @Accept json
// @Produce json
// @Param input body dto.CreateAdminInput true "Данные администратора"
// @Success 201 {object} dto.CreateAdminOutput
// @Failure 400 {object} dto.ErrorOutput
// @Failure 500 {object} dto.ErrorOutput
// @Router /secret/create-admin [post]
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
