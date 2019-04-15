package middleware

import (
	"github.com/gin-gonic/gin"
)

// RecoveryHandlerFunc is called if the recovery finds an error
// This should write a default response
type RecoveryHandlerFunc func(*gin.Context, interface{}) error

// Recovery handles a panic and returns a 500 error
func Recovery(fn RecoveryHandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func(c *gin.Context, fn RecoveryHandlerFunc) {
			err := recover()
			if err == nil {
				return
			}

			err = fn(c, err)
			if err != nil {
				panic(err)
			}
		}(c, fn)

		c.Next()
	}
}
