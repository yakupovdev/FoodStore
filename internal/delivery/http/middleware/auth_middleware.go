package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yakupovdev/FoodStore/internal/domain/entity"
	"github.com/yakupovdev/FoodStore/internal/domain/service"
	"github.com/yakupovdev/FoodStore/internal/usecase"
)

func AuthMiddleware(validator usecase.TokenValidator, tokenSvc service.TokenService, tokenType entity.TokenType) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing Authorization header",
			})
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid Authorization format",
			})
			return
		}

		claims, err := tokenSvc.ParseToken(parts[1], tokenType)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			return
		}

		if tokenType == entity.AccessTokenType {
			isValid, err := validator.IsTokenValid(context.Background(), claims.UserID, parts[1])
			if err != nil || !isValid {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "invalid or revoked token",
				})
				return
			}
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_type", claims.UserType)

		c.Next()
	}
}
