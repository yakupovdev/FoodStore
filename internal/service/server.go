package service

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yakupovdev/FoodStore/internal/handler/auth_handler"
	"github.com/yakupovdev/FoodStore/internal/handler/mail_handler"
	"github.com/yakupovdev/FoodStore/internal/handler/security_handler"
	"github.com/yakupovdev/FoodStore/internal/repository"
)

type Server struct {
	srv *http.Server
}

func SetupRouter(pg *repository.Postgres) *gin.Engine {
	router := gin.Default()

	auth := router.Group("/auth")
	{
		auth.POST("/register", func(c *gin.Context) {
			auth_handler.RegisterHandlers(c, pg)
		})

		auth.POST("/login", func(c *gin.Context) {
			auth_handler.LoginHandlers(c, pg)
		})

	}
	protected := router.Group("/protected")
	protected.Use(AuthMiddleware())
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
		recovery.POST("/send-code", func(c *gin.Context) {
			mail_handler.SendCode(c, pg)
		})
		recovery.POST("/verify-code", func(c *gin.Context) {
			mail_handler.VerifyCode(c, pg)

		})
	}
	recovery_password := router.Group("/recovery-password")
	recovery_password.Use(AuthMiddleware())
	{
		recovery_password.POST("/reset", func(c *gin.Context) {
			security_handler.ResetPassword(c, pg)
		})
	}

	return router
}

func NewServer(addr string, router *gin.Engine) *Server {
	if router == nil {
		router = gin.Default()
	}
	return &Server{
		srv: &http.Server{
			Addr:         addr,
			Handler:      router,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}

}

func (s *Server) Run() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
