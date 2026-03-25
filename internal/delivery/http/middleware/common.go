package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yakupovdev/FoodStore/internal/domain/logger"
	"go.uber.org/zap"
)

const requestIDHeader = "X-Request-ID"

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(requestIDHeader)
		if requestID == "" {
			requestID = uuid.NewString()
		}
		c.Request.Header.Set(requestIDHeader, requestID)
		c.Writer.Header().Set(requestIDHeader, requestID)
		c.Next()
	}
}

func Logger(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(requestIDHeader)
		l := log.With(
			zap.String("request_id", requestID),
			zap.String("url", c.Request.URL.String()),
			zap.String("method", c.Request.Method),
			zap.String("client_ip", c.ClientIP()),
		)
		ctx := context.WithValue(c.Request.Context(), "log", l)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func Trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		log := logger.FromContext(ctx)

		before := time.Now()
		log.Info(">>> incoming HTTP request",
			zap.Time("time", before.UTC()),
		)

		c.Next()

		after := time.Now()
		log.Info("<<< completed HTTP request",
			zap.Time("time", after.UTC()),
			zap.Duration("duration", after.Sub(before)),
			zap.Int("status", c.Writer.Status()),
		)
	}
}
