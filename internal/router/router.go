package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yakupovdev/FoodStore/internal/controller"
	"github.com/yakupovdev/FoodStore/internal/middleware"
	"github.com/yakupovdev/FoodStore/internal/security"
	"github.com/yakupovdev/FoodStore/internal/usecase"
)

type Deps struct {
	AuthController               *controller.AuthController
	EmailController              *controller.EmailController
	RecoveryController           *controller.RecoveryController
	RefreshAccessTokenController *controller.RefreshAccessTokenController
	CheckTokenIsValidUsecase     *usecase.CheckTokenIsValidUsecase
	ClientController             *controller.ClientController
}

func SetupRouter(d Deps) *gin.Engine {
	router := gin.Default()

	auth := router.Group("/auth")
	{
		auth.POST("/register", d.AuthController.RegisterUser)

		auth.POST("/login", d.AuthController.LoginUser)

	}
	protected := router.Group("/protected")
	protected.Use(middleware.AuthMiddleware(d.CheckTokenIsValidUsecase, security.AccessToken))
	{
		protected.GET("/profile", func(c *gin.Context) {
			userID, exists := c.Get("user_id")
			if !exists {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id not found in context"})
				return
			}
			userType, exists := c.Get("user_type")
			if !exists {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "user_type not found in context"})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"user_id":   userID,
				"user_type": userType,
				"message":   "This is a protected profile route",
			})
		})
	}
	client := router.Group("/client")
	client.Use(middleware.AuthMiddleware(d.CheckTokenIsValidUsecase, security.AccessToken))
	{
		client.GET("/orders", d.ClientController.GetOrders)
		client.GET("/profile", d.ClientController.GetProfile)
	}

	recovery := router.Group("/recovery")
	{
		recovery.POST("/send-code", d.EmailController.SendCodeByEmail)
		recovery.POST("/verify-code", d.EmailController.VerifyCode)
	}
	recovery_password := router.Group("/recovery-password")
	recovery_password.Use(middleware.AuthMiddleware(d.CheckTokenIsValidUsecase, security.RecoveryToken))
	{
		recovery_password.POST("/reset", d.RecoveryController.ResetUserPassword)
	}

	refreshAccess := router.Group("/refresh-access")
	refreshAccess.Use(middleware.AuthMiddleware(d.CheckTokenIsValidUsecase, security.RefreshToken))
	{
		refreshAccess.POST("/token", d.RefreshAccessTokenController.RefreshAccessToken)
	}

	return router
}
