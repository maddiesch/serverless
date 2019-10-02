package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	// AccessControlAllowMethodsDefaults contains all the default HTTP methods allowed for CORS
	AccessControlAllowMethodsDefaults = []string{"GET", "POST", "PATCH", "OPTIONS", "DELETE"}

	// AccessControlAllowHeadersDefaults is the allowed headers for CORS
	AccessControlAllowHeadersDefaults = []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"}

	// DefaultCorsConfiguration is the default configuration for CORS middleware
	DefaultCorsConfiguration = &CorsConfiguration{
		AccessControlAllowOrigin:  "*",
		AccessControlAllowMethods: AccessControlAllowMethodsDefaults,
		AccessControlAllowHeaders: AccessControlAllowHeadersDefaults,
	}
)

var (
	// CORSOriginValidator can be used to validate if the origin is valid
	CORSOriginValidator func(string) bool
)

const (
	accessControlAllowOriginHeader  = "Access-Control-Allow-Origin"
	accessControlAllowMethodsHeader = "Access-Control-Allow-Methods"
	accessControlAllowHeadersHeader = "Access-Control-Allow-Headers"
)

// CorsConfiguration contains all the options for CORS middleware
type CorsConfiguration struct {
	AccessControlAllowOrigin  string
	AccessControlAllowMethods []string
	AccessControlAllowHeaders []string
}

// Cors middleware handles adding headers for requests.
func Cors(config *CorsConfiguration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(config.AccessControlAllowOrigin) > 0 {
			c.Header(accessControlAllowOriginHeader, config.AccessControlAllowOrigin)
		}
		if CORSOriginValidator != nil {
			origin := c.GetHeader("Origin")
			if origin != "" && CORSOriginValidator(origin) {
				c.Header(accessControlAllowOriginHeader, origin)
			}
		}
		if len(config.AccessControlAllowMethods) > 0 {
			c.Header(accessControlAllowMethodsHeader, strings.Join(config.AccessControlAllowMethods, ", "))
		}
		if len(config.AccessControlAllowHeaders) > 0 {
			c.Header(accessControlAllowHeadersHeader, strings.Join(config.AccessControlAllowHeaders, ", "))
		}

		c.Next()
	}
}
