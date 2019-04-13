package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// Runtime adds a `X-Runtime` header to the response
func Runtime() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		defer func(start time.Time) {
			dur := time.Now().Sub(start)
			c.Header("X-Runtime", fmt.Sprintf("%0.02fms", dur.Seconds()*1e3))
		}(start)

		c.Next()
	}
}
