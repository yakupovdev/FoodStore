package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func extractUserID(ctx *gin.Context) (int64, bool) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		return 0, false
	}

	switch v := userIDVal.(type) {
	case int64:
		return v, true
	case int:
		return int64(v), true
	case float64:
		return int64(v), true
	case string:
		parsed, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, false
		}
		return parsed, true
	default:
		return 0, false
	}
}

func extractUserType(ctx *gin.Context) (string, bool) {
	userTypeVal, exists := ctx.Get("user_type")
	if !exists {
		return "", false
	}

	utype, ok := userTypeVal.(string)
	return utype, ok
}
