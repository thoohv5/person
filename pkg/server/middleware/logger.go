package middleware

import (
	"bytes"
	"net/http/httputil"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/thoohv5/person/pkg/log"
)

type Fn func(c *gin.Context) []log.IField

type Config struct {
	TimeFormat string
	UTC        bool
	SkipPaths  []string
	Context    Fn
}

func Logger(logger log.Logger, timeFormat string, utc bool) gin.HandlerFunc {
	return LoggerWithConfig(logger, &Config{TimeFormat: timeFormat, UTC: utc})
}

func LoggerWithConfig(logger log.Logger, conf *Config) gin.HandlerFunc {
	skipPaths := make(map[string]bool, len(conf.SkipPaths))
	for _, path := range conf.SkipPaths {
		skipPaths[path] = true
	}

	return func(c *gin.Context) {
		// some evil middlewares modify this values
		path := c.Request.URL.Path

		requestBr, _ := httputil.DumpRequest(c.Request, true)
		c.Next()
		var responseBr []byte
		if resp := c.Request.Response; resp != nil {
			rb, _ := httputil.DumpResponse(resp, false)
			responseBr = rb
		}

		if _, ok := skipPaths[path]; !ok {
			end := time.Now()
			if conf.UTC {
				end = end.UTC()
			}

			fields := []log.Field{
				func(field log.IField) {
					field.Set("request", string(bytes.ReplaceAll(requestBr, []byte("\n"), []byte(" "))), log.WithType(reflect.String))
					field.Set("response", string(bytes.ReplaceAll(responseBr, []byte("\n"), []byte(" "))), log.WithType(reflect.String))
				},
			}
			if conf.TimeFormat != "" {
				fields = append(fields, func(field log.IField) {
					field.Set("time", end.Format(conf.TimeFormat))
				})
			}

			if conf.Context != nil {
				fields = append(fields, func(field log.IField) {
					field.Set("", c)
				})
			}

			if len(c.Errors) > 0 {
				// Append error field if this is an erroneous request.
				for _, e := range c.Errors.Errors() {
					logger.Errorc(c, e, fields...)
				}
			} else {
				logger.Infoc(c, path, fields...)
			}
		}
	}
}
