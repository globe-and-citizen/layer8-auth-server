package ginUtils

import (
	"context"
	"fmt"
	"globe-and-citizen/layer8/auth-server/pkg/log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
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

func AccessLogJSON(p gin.LogFormatterParams) string {
	return fmt.Sprintf(`{"time":"%s", "status":%d, "latency":"%s", "ip":"%s", "method":"%s", "path":"%s", "body_size":"%d", "request_id":"%s"}%s`,
		p.TimeStamp.Format(time.RFC3339),
		p.StatusCode,
		p.Latency,
		p.ClientIP,
		p.Method,
		p.Path,
		p.BodySize,
		p.Keys[RequestIDKey],
		"\n",
	)
}

func AccessLogText(p gin.LogFormatterParams) string {
	return fmt.Sprintf("[%s] %d | %13v | %15s | body_size=%13d | req_id=%s | %-7s %s\n",
		p.TimeStamp.Format(time.RFC3339),
		p.StatusCode,
		p.Latency,
		p.ClientIP,
		p.BodySize,
		p.Keys[RequestIDKey],
		p.Method,
		p.Path,
	)
}

func AccessLog(l log.ILogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next() // handle request

		// Skip logging for internal/health endpoints
		if path == "/health" || path == "/config.js" || strings.HasPrefix(path, "/assets") {
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
			log.F("request_id", c.GetString(RequestIDKey)),
		}

		// Log errors and slow requests; skip successful fast requests
		isError := status >= 400
		isSlow := latency > 500*time.Millisecond

		if isError {
			fields = append(fields,
				log.F("ip", c.ClientIP()),
				log.F("user_agent", c.Request.UserAgent()),
			)
			if ct := c.GetHeader("Content-Type"); ct != "" {
				fields = append(fields, log.F("content_type", ct))
			}
			if at := getAuthorizationType(c.GetHeader("Authorization")); at != "" {
				fields = append(fields, log.F("auth_type", at))
			}

			// Log error details
			if len(c.Errors) > 0 {
				l.Error("http request failed", c.Errors.Last(), fields...)
			} else {
				l.Warn("http request error", fields...)
			}
			return
		}

		if isSlow {
			fields = append(fields,
				log.F("ip", c.ClientIP()),
				log.F("response_size", c.Writer.Size()),
			)
			l.Warn("http request slow", fields...)
		}
	}
}

// OTel is a middleware that traces HTTP requests and measures latency using OpenTelemetry
func OTel(tracer trace.Tracer) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the remote parent from traceparent/tracestate.
		ctx := otel.GetTextMapPropagator().Extract(
			c.Request.Context(),
			propagation.HeaderCarrier(c.Request.Header),
		)

		// Create the server span as a child of the extracted context.
		ctx, span := tracer.Start(
			ctx,
			fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.Path),
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				attribute.String("http.method", c.Request.Method),
				attribute.String("http.url", c.Request.URL.String()),
				attribute.String("http.target", c.Request.RequestURI),
				attribute.String("http.host", c.Request.Host),
				attribute.String("http.scheme", c.Request.URL.Scheme),
				attribute.String("http.client_ip", c.ClientIP()),
				attribute.String("http.user_agent", c.Request.UserAgent()),
			),
		)
		defer span.End()

		// Update the request context with the span
		c.Request = c.Request.WithContext(ctx)

		// Record the start time for latency measurement
		start := time.Now()

		// Call the next handler
		c.Next()

		// Record response details
		latency := time.Since(start)
		statusCode := c.Writer.Status()

		span.SetAttributes(
			attribute.Int("http.status_code", statusCode),
			attribute.Int64(
				"http.response_content_length",
				int64(c.Writer.Size()),
			),
		)

		if statusCode >= 400 {
			span.RecordError(fmt.Errorf("HTTP %d", statusCode))
			span.SetStatus(codes.Error, fmt.Sprintf("HTTP %d", statusCode))
		} else {
			span.SetStatus(codes.Ok, "OK")
		}

		// Store latency in context for access log middleware
		c.Set("latency", latency)
	}
}
