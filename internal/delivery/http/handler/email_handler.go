package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yakupovdev/FoodStore/internal/delivery/http/dto"
	"github.com/yakupovdev/FoodStore/internal/domain"
	"github.com/yakupovdev/FoodStore/internal/usecase"
)

type EmailHandler struct {
	uc *usecase.RecoveryUsecase
}

func NewEmailHandler(uc *usecase.RecoveryUsecase) *EmailHandler {
	return &EmailHandler{uc: uc}
}

// SendCodeByEmail godoc
// @Summary Отправить код восстановления
// @Description Отправляет код восстановления пароля на email пользователя
// @Tags Recovery
// @Accept json
// @Produce json
// @Param input body dto.SendCodeInput true "Данные для отправки кода"
// @Success 200 {object} dto.SendCodeOutput
// @Failure 400 {object} dto.ErrorOutput "Пустые поля"
// @Failure 404 {object} dto.ErrorOutput "Пользователь не найден"
// @Failure 500 {object} dto.ErrorOutput
// @Router /recovery/send-code [post]
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

// VerifyCode godoc
// @Summary Проверить код восстановления
// @Description Проверяет код восстановления и возвращает recovery токен
// @Tags Recovery
// @Accept json
// @Produce json
// @Param input body dto.VerifyCodeInput true "Данные для проверки кода"
// @Success 200 {object} dto.VerifyCodeOutput
// @Failure 400 {object} dto.ErrorOutput "Пустые поля или неверный код"
// @Failure 404 {object} dto.ErrorOutput "Пользователь не найден"
// @Failure 500 {object} dto.ErrorOutput
// @Router /recovery/verify-code [post]
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
