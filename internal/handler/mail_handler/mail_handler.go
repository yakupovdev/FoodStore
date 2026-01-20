package mail_handler

import (
	"fmt"
	"net/http"
	"net/smtp"

	"github.com/gin-gonic/gin"
	"github.com/yakupovdev/FoodStore/internal/model"
	"github.com/yakupovdev/FoodStore/internal/repository"
	"github.com/yakupovdev/FoodStore/internal/security"
	"github.com/yakupovdev/FoodStore/internal/storage"
)

func SendCode(c *gin.Context, pg *repository.Postgres) {
	var req model.VerifyEmailRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	code := security.GenerateAccessCodeByEmail()
	to := req.Email
	auth := smtp.PlainAuth(
		"",
		"foodstorewwgo@gmail.com",
		"dkeywmbvieuiuazj",
		"smtp.gmail.com",
	)

	msg := []byte(
		"From: FoodStore <foodstorewwgo@gmail.com>\r\n" +
			"To:" + to + "\r\n" +
			"Subject: Test Email from FoodStore\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-UserType: text/plain; charset=\"UTF-8\"\r\n" +
			"\r\n" +
			"Here your code: " + code + "\n",
	)

	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		"foodstorewwgo@gmail.com",
		[]string{to},
		msg,
	)
	if err != nil {
		fmt.Println(err)
	}
	userID, err := pg.GetUserIDByEmail(to)
	if err != nil {
		return
	}
	err = pg.SaveRecoveryCode(userID, to, code)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func VerifyCode(c *gin.Context, pg *repository.Postgres) {
	var req model.VerifyCodeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, err := pg.GetUserIDByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": storage.ErrDatabaseConnection.Error(),
		})
		return
	}

	isValidCode, err := pg.VerifyRecoveryCode(req.Email, req.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": storage.ErrDatabaseConnection.Error(),
		})
		return
	}
	if !isValidCode {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": security.ErrInvalidRecoveryCode.Error(),
		})
		return
	}

	token, err := security.GenerateToken(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": security.ErrTokenGeneration.Error(),
		})
		return
	}

	err = pg.DeleteExpiredRecoveryCodes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": storage.ErrDatabaseConnection.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
