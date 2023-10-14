package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
)

// Transform 转换中间件
func Transform(callback func(ctx context.Context, data map[string][]string) context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request = c.Request.WithContext(callback(c.Request.Context(), c.Request.Header))
		c.Next()
	}
}
