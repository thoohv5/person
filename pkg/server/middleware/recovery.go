package middleware

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/thoohv5/person/pkg/log"
)

func defaultHandleRecovery(c *gin.Context, err interface{}) {
	c.AbortWithStatus(http.StatusInternalServerError)
}

func Recovery(logger log.Logger, stack bool) gin.HandlerFunc {
	return CustomRecovery(logger, stack, defaultHandleRecovery)
}

func CustomRecovery(logger log.Logger, stack bool, recovery gin.RecoveryFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.Errorc(c, c.Request.URL.Path,
						func(fields log.IField) {
							fields.Set("error", err)
							fields.Set("request", string(httpRequest))
						},
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					logger.Errorc(c, "[Recovery from panic]",
						func(fields log.IField) {
							fields.Set("time", time.Now())
							fields.Set("error", err)
							fields.Set("request", string(httpRequest))
							fields.Set("stack", string(debug.Stack()))
						},
					)
				} else {
					logger.Errorc(c, "[Recovery from panic]",
						func(fields log.IField) {
							fields.Set("time", time.Now())
							fields.Set("error", err)
							fields.Set("request", string(httpRequest))
						},
					)
				}
				recovery(c, err)
			}
		}()
		c.Next()
	}
}
