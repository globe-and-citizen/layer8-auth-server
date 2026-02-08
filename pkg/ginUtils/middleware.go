package ginUtils

import (
	"context"
	"globe-and-citizen/layer8/auth-server/pkg/log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const RequestIDKey = "request_id"

func RequestID(c *gin.Context) {
	reqID := getRequestID(c)

	// Store in gin.Context
	c.Set(RequestIDKey, reqID)

	// Store in request context
	ctx := context.WithValue(c.Request.Context(), RequestIDKey, reqID)
	c.Request = c.Request.WithContext(ctx)

	// Return to client
	c.Header("X-Request-ID", reqID)

	c.Next()
}

func AccessLog(l log.ILogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next() // handle request

		if path == "/health" || strings.HasPrefix(path, "assets/") {
			return
		}

		latency := time.Since(start)
		status := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		fields := []log.Field{
			log.F("status", status),
			log.F("method", c.Request.Method),
			log.F("path", path),
			log.F("latency_ms", latency.Milliseconds()),
			log.F("ip", c.ClientIP()),
			log.F("user_agent", c.Request.UserAgent()),
			log.F("request_id", c.GetString(RequestIDKey)),
		}

		ct := c.GetHeader("Content-Type")
		if ct != "" {
			fields = append(fields, log.F("content_type", ct))
		}

		if at := getAuthorizationType(c); at != "" {
			fields = append(fields, log.F("auth_type", at))
		}

		// Gin errors (if any)
		if len(c.Errors) > 0 {
			l.Error(
				"http request failed",
				c.Errors.Last(),
				fields...,
			)
			return
		}

		l.Info("access_log", fields...)
	}
}
