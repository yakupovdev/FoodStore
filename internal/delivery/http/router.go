package http

import (
	"github.com/gin-gonic/gin"
	"github.com/yakupovdev/FoodStore/internal/delivery/http/handler"
	"github.com/yakupovdev/FoodStore/internal/delivery/http/middleware"
	"github.com/yakupovdev/FoodStore/internal/domain/entity"
	"github.com/yakupovdev/FoodStore/internal/domain/service"
	"github.com/yakupovdev/FoodStore/internal/usecase"
)

type RouterDeps struct {
	AuthHandler         *handler.AuthHandler
	EmailHandler        *handler.EmailHandler
	RecoveryHandler     *handler.RecoveryHandler
	RefreshTokenHandler *handler.RefreshTokenHandler
	ClientHandler       *handler.ClientHandler
	SellerHandler       *handler.SellerHandler
	AdminHandler        *handler.AdminHandler
	TokenValidator      usecase.TokenValidator
	TokenService        service.TokenService
}

func SetupRouter(d RouterDeps) *gin.Engine {
	router := gin.Default()

	auth := router.Group("/auth")
	{
		auth.POST("/register", d.AuthHandler.RegisterUser)
		auth.POST("/login", d.AuthHandler.LoginUser)
	}

	protected := router.Group("/protected")
	protected.Use(middleware.AuthMiddleware(d.TokenValidator, d.TokenService, entity.AccessTokenType))
	{
		client := protected.Group("/client")
		client.Use(middleware.AccessTypeMiddleware("client"))
		{
			client.GET("/orders", d.ClientHandler.GetOrders)
			client.POST("/orders", d.ClientHandler.CreateOrder)
			client.GET("/profile", d.ClientHandler.GetProfile)
			client.GET("/products", d.ClientHandler.GetProducts)
			client.POST("/address", d.ClientHandler.AddAdress)
			client.POST("/cart", d.ClientHandler.AddToCart)
			client.GET("/cart", d.ClientHandler.GetCartItems)
			client.GET("/products/:product_id", d.ClientHandler.GetProductByID)
			client.GET("/", d.ClientHandler.GetProductsByPriority)
		}

		seller := protected.Group("/seller")
		seller.Use(middleware.AccessTypeMiddleware("seller"))
		{
			seller.GET("/profile", d.SellerHandler.GetProfile)
			seller.GET("/offers", d.SellerHandler.GetOffers)
			seller.POST("/offers", d.SellerHandler.CreateOfferByExistProducts)
			seller.PUT("/offers", d.SellerHandler.UpdateOffer)
			seller.DELETE("/offers", d.SellerHandler.DeleteOffer)
			seller.GET("/products", d.SellerHandler.GetExistProducts)
			seller.POST("/products", d.SellerHandler.CreateOffer)
			seller.GET("/subscription", d.SellerHandler.PurchaseSubscription)
		}

		admin := protected.Group("/admin")
		admin.Use(middleware.AccessTypeMiddleware("admin"))
		{
			admin.DELETE("/users", d.AdminHandler.DeleteUser)
			admin.POST("/balance", d.AdminHandler.UpdateBalance)
			admin.GET("/logs", d.AdminHandler.GetAllLogTransactions)
			admin.GET("/users", d.AdminHandler.GetAllUsers)
		}
	}

	secret := router.Group("/secret")
	{
		secret.POST("/create-admin", d.AdminHandler.CreateAdmin)
	}

	recovery := router.Group("/recovery")
	{
		recovery.POST("/send-code", d.EmailHandler.SendCodeByEmail)
		recovery.POST("/verify-code", d.EmailHandler.VerifyCode)
	}

	recoveryPassword := router.Group("/recovery-password")
	recoveryPassword.Use(middleware.AuthMiddleware(d.TokenValidator, d.TokenService, entity.RecoveryTokenType))
	{
		recoveryPassword.POST("/reset", d.RecoveryHandler.ResetUserPassword)
	}

	refreshAccess := router.Group("/refresh-access")
	refreshAccess.Use(middleware.AuthMiddleware(d.TokenValidator, d.TokenService, entity.RefreshTokenType))
	{
		refreshAccess.POST("/token", d.RefreshTokenHandler.RefreshAccessToken)
	}

	return router
}
