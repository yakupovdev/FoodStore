package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AccessTypeMiddleware(userType string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uType, exist := ctx.Get("user_type")

		if !exist {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"Message": "User type not exist",
			})
			ctx.Abort()
			return
		}

		if uType != userType {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"Message": "User type not match",
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
