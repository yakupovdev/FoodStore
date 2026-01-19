package auth_handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yakupovdev/FoodStore/internal/model"
	"github.com/yakupovdev/FoodStore/internal/repository"
	"github.com/yakupovdev/FoodStore/internal/security"
	"github.com/yakupovdev/FoodStore/internal/service"
	"github.com/yakupovdev/FoodStore/internal/storage"
)

func RegisterHandlers(c *gin.Context, pg *repository.Postgres) {
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
	exist, err := pg.GetUserByEmail(req.Email)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if exist {
		c.JSON(http.StatusConflict, gin.H{
			"error": repository.ErrDuplicateLogin.Error(),
		})
		return
	}

	_, err = pg.CreateUser(req.Email, hashHex, req.Type, req.Balance)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered"})
}

func LoginHandlers(c *gin.Context, pg *repository.Postgres) {
	var req model.AuthRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashHex := security.HashPassword(req.Password)
	if pg == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": storage.ErrDatabaseConnection.Error(),
		})
		return
	}
	log.Println(hashHex)
	userID, err := pg.LoginUser(req.Email, hashHex)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": service.ErrInvalidCredentials.Error(),
		})
		return
	}

	token, err := security.GenerateToken(int64(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": security.ErrTokenGeneration.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": token,
		"message":      "User logged in",
	})
}
