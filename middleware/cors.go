package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	AccessControlAllowMethodsDefaults = []string{"GET", "POST", "PATCH", "OPTIONS", "DELETE"}

	AccessControlAllowHeadersDefaults = []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"}

	DefaultCorsConfiguration = &CorsConfiguration{
		AccessControlAllowOrigin:  "*",
		AccessControlAllowMethods: AccessControlAllowMethodsDefaults,
		AccessControlAllowHeaders: AccessControlAllowHeadersDefaults,
	}
)

type CorsConfiguration struct {
	AccessControlAllowOrigin  string
	AccessControlAllowMethods []string
	AccessControlAllowHeaders []string
}

func Cors(config *CorsConfiguration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(config.AccessControlAllowOrigin) > 0 {
			c.Header("Access-Control-Allow-Origin", config.AccessControlAllowOrigin)
		}
		if len(config.AccessControlAllowMethods) > 0 {
			c.Header("Access-Control-Allow-Methods", strings.Join(config.AccessControlAllowMethods, ", "))
		}
		if len(config.AccessControlAllowHeaders) > 0 {
			c.Header("Access-Control-Allow-Headers", strings.Join(config.AccessControlAllowHeaders, ", "))
		}

		c.Next()
	}
}
