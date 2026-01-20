package security_handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yakupovdev/FoodStore/internal/model"
	"github.com/yakupovdev/FoodStore/internal/repository"
	"github.com/yakupovdev/FoodStore/internal/security"
)

func ResetPassword(c *gin.Context, pg *repository.Postgres) {
	var req model.ResetUserPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(500, gin.H{"error": "user_id not found in context"})
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
			c.JSON(500, gin.H{"error": "invalid user_id format"})
			return
		}
		uid = parsed
	default:
		c.JSON(500, gin.H{"error": "unsupported user_id type"})
		return
	}

	hashPassword := security.HashPassword(req.NewPassword)

	err := pg.UpdateUserPassword(uid, hashPassword)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Password reset successful"})
}
