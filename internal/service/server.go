package service

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yakupovdev/FoodStore/internal/handler"
	"github.com/yakupovdev/FoodStore/internal/storage"
)

type Server struct {
	srv *http.Server
}

func SetupRouter(pg *storage.Postgres) *gin.Engine {
	router := gin.Default()

	auth := router.Group("/auth")
	{
		auth.POST("/register", func(c *gin.Context) {
			handler.RegisterHandlers(c, pg)
		})

		auth.POST("/login", func(c *gin.Context) {
			handler.LoginHandlers(c, pg)
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
