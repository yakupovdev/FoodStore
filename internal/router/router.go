package router

import (
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
	SellerController             *controller.SellerController
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
		client := protected.Group("/client")
		client.Use(middleware.AccessTypeMiddleware("client"))
		{
			client.GET("/orders", d.ClientController.GetOrders)
			client.GET("/profile", d.ClientController.GetProfile)
		}

		seller := protected.Group("/seller")
		seller.Use(middleware.AccessTypeMiddleware("seller"))
		{
			seller.GET("/profile", d.SellerController.GetProfile)
			seller.GET("/offers", d.SellerController.GetOffers)
			seller.POST("/offers", d.SellerController.CreateOffer)
		}
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
