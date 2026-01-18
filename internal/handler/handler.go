package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yakupovdev/FoodStore/internal/model"
	"github.com/yakupovdev/FoodStore/internal/security"
	"github.com/yakupovdev/FoodStore/internal/storage"
)

func RegisterHandlers(c *gin.Context, pg *storage.Postgres) {
	var req model.UserData

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hashHex := security.HashPassword(req.Password)
	if pg == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database not initialized"})
		return
	}
	log.Println("vatahell")
	exist, err := pg.GetUserByLogin(req.Login)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if exist {
		c.JSON(http.StatusConflict, gin.H{"error": model.ErrDuplicateLogin})
		return
	}

	_, err = pg.CreateUser(req.Login, hashHex, req.Type, req.Balance)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered"})
}

func LoginHandlers(c *gin.Context, pg *storage.Postgres) {
	var req model.AuthRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hashHex := security.HashPassword(req.Password)
	if pg == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": model.DatabaseConnectionError})
		return
	}
	log.Println(hashHex)
	userID, err := pg.LoginUser(req.Login, hashHex)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": model.ErrInvalidCredentials})
		return
	}

	token, err := security.GenerateToken(int64(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": model.CouldNotGenerateTokenError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": token,
		"message":      "User logged in",
	})
}
