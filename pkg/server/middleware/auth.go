package middleware

import (
	"context"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// Auth 鉴权
func Auth(opts ...Option) gin.HandlerFunc {

	o := &options{
		skipUrls: make([]string, 0),
	}
	for _, opt := range opts {
		opt(o)
	}

	mSkipUrls := make(map[string]bool)
	for _, url := range o.skipUrls {
		mSkipUrls[url] = true
	}

	return func(c *gin.Context) {
		for _, url := range o.skipUrls {
			if ok, err := filepath.Match(url, c.Request.URL.Path); ok && err == nil {
				c.Next()
				return
			}
		}

		if getClientIP := o.getClientIP; getClientIP != nil {
			if clientIP := getClientIP(c.Request.Context()); clientIP == "127.0.0.1" {
				c.Next()
				return
			}
		}

		cUserID := ""
		getUserID := o.getUserID
		if getUserID == nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		cUserID = getUserID(c.Request.Context())

		// 如果请求中没有 UserID，则返回状态码 401 Unauthorized
		if len(cUserID) == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if cb := o.callback; cb == nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		_, boolean := o.callback(c.Request.Context(), cUserID)
		if !boolean {
			// 检验失败
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// 用户允许登录，放行请求
		c.Next()
	}
}

type options struct {
	getClientIP func(c context.Context) string
	getUserID   func(c context.Context) string
	callback    func(ctx context.Context, userID string) (context.Context, bool)
	skipUrls    []string
}

type Option func(*options)

func WithUserID(getUserID func(c context.Context) string) Option {
	return func(o *options) {
		o.getUserID = getUserID
	}
}

func WithClientIP(getClientIP func(c context.Context) string) Option {
	return func(o *options) {
		o.getClientIP = getClientIP
	}
}

func WithCallback(callback func(ctx context.Context, userID string) (context.Context, bool)) Option {
	return func(o *options) {
		o.callback = callback
	}
}

func WithSkipUrls(skipUrls []string) Option {
	return func(o *options) {
		o.skipUrls = skipUrls
	}
}
