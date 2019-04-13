package middleware

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger logs the request
func Logger(fn func(string)) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		requestID := c.GetString(GinContextRequestIDKey)
		if len(requestID) == 0 {
			requestID = "{missing-request-id}"
		}

		msg := fmt.Sprintf("%s [%s] %s%s", requestID, c.Request.Method, c.Request.URL.Path, c.Request.URL.RawQuery)

		fn(msg)

		defer func(c *gin.Context, requestID string, start time.Time) {
			if strings.Compare(gin.Mode(), gin.ReleaseMode) == 0 {
				// return
			}
			dur := time.Now().Sub(start)
			msg := fmt.Sprintf("%s Completed in %sms [%d]", requestID, fmt.Sprintf("%0.03f", dur.Seconds()*1e3), c.Writer.Status())

			fn(msg)
		}(c, requestID, start)

		c.Next()
	}
}
