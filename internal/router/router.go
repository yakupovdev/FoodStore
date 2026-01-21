package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yakupovdev/FoodStore/internal/controller"
	middleware2 "github.com/yakupovdev/FoodStore/internal/middleware"
)

type Deps struct {
	AuthController     *controller.AuthController
	EmailController    *controller.EmailController
	RecoveryController *controller.RecoveryController
}

func SetupRouter(d Deps) *gin.Engine {
	router := gin.New()

	router.Use(middleware2.LoggerMiddleware())

	auth := router.Group("/auth")
	{
		auth.POST("/register", d.AuthController.RegisterUser)

		auth.POST("/login", d.AuthController.LoginUser)

	}
	protected := router.Group("/protected")
	protected.Use(middleware2.AuthMiddleware())
	{
		protected.GET("/profile", func(c *gin.Context) {
			userID, exists := c.Get("user_id")
			if !exists {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id not found in context"})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"user_id": userID,
				"message": "This is a protected profile route",
			})
		})
	}
	recovery := router.Group("/recovery")
	{
		recovery.POST("/send-code", d.EmailController.SendCodeByEmail)
		recovery.POST("/verify-code", d.EmailController.VerifyCode)
	}
	recovery_password := router.Group("/recovery-password")
	recovery_password.Use(middleware2.AuthMiddleware())
	{
		recovery_password.POST("/reset", d.RecoveryController.ResetUserPassword)
	}

	return router
}
