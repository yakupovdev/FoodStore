package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yakupovdev/FoodStore/internal/model"
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
			var req model.UserData
			hash := sha256.Sum256([]byte(req.Password))
			hashHex := hex.EncodeToString(hash[:])
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			if pg == nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database not initialized"})
				return
			}
			_, err := pg.Create(req.Login, hashHex, req.Type, req.Balance)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": "User registered"})
		})

		auth.POST("/login", func(c *gin.Context) {
			var req model.AuthRequest
			hash := sha256.Sum256([]byte(req.Password))
			hashHex := hex.EncodeToString(hash[:])
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			if pg == nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database not initialized"})
				return
			}
			userID, err := pg.Login(req.Login, hashHex)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
				return
			}

			token, err := GenerateToken(int64(userID))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "could not generate token",
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"access_token": token,
			})

			c.JSON(http.StatusOK, gin.H{"message": "User logged in"})
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
