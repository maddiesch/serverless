package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"
)

const (
	// GinContextRequestIDKey Key for the request id on the context
	GinContextRequestIDKey = "serverless.request-id"
)

// RequestID adds a random ID to request.
// It is available on the gin context using the `GinContextRequestIDKey`
// it is also returned in the `X-Request-Id` HTTP header.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := ksuid.New().String()

		c.Set(GinContextRequestIDKey, requestID)
		c.Header("X-Request-Id", requestID)

		c.Next()
	}
}
