package ginUtils

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func getAuthorizationType(authHeader string) string {
	if authHeader == "" {
		return ""
	}

	// Split at first space: "<type> <credentials>"
	i := strings.IndexByte(authHeader, ' ')
	if i == -1 {
		// Malformed but still indicates scheme
		return strings.ToLower(authHeader)
	}

	return strings.ToLower(authHeader[:i])
}

func getRequestID(c *gin.Context) string {
	// 1️⃣ Client / proxy provided
	if id := c.GetHeader("X-Request-ID"); id != "" {
		return id
	}

	if id := c.GetHeader("X-Correlation-ID"); id != "" {
		return id
	}

	// 2️⃣ Trace context (W3C)
	if tp := c.GetHeader("traceparent"); tp != "" {
		if traceID := extractTraceID(tp); traceID != "" {
			return traceID
		}
	}

	// 3️⃣ Generate
	return uuid.NewString()
}

func extractTraceID(traceparent string) string {
	// Format:
	// version-traceid-spanid-flags
	// 00-4bf92f3577b34da6a3ce929d0e0e4736-...

	parts := strings.Split(traceparent, "-")
	if len(parts) < 2 {
		return ""
	}

	traceID := parts[1]
	if len(traceID) == 32 {
		return traceID
	}

	return ""
}
