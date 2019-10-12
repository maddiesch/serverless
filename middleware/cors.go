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
)

// CORSOriginProvider can be used to validate if the origin is valid
type CORSOriginProvider func(string) string

const (
	accessControlAllowOriginHeader  = "Access-Control-Allow-Origin"
	accessControlAllowMethodsHeader = "Access-Control-Allow-Methods"
	accessControlAllowHeadersHeader = "Access-Control-Allow-Headers"
)

// CORSConfiguration contains all the options for CORS middleware
type CORSConfiguration struct {
	OriginProvider            CORSOriginProvider
	AccessControlAllowMethods []string
	AccessControlAllowHeaders []string
}

// DefaultCORS returns a default CORS configuration
func DefaultCORS(fn CORSOriginProvider) gin.HandlerFunc {
	return CORS(&CORSConfiguration{
		OriginProvider:            fn,
		AccessControlAllowMethods: AccessControlAllowMethodsDefaults,
		AccessControlAllowHeaders: AccessControlAllowHeadersDefaults,
	})
}

// CORS middleware handles adding headers for requests.
func CORS(config *CORSConfiguration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if origin := config.OriginProvider(c.GetHeader("Origin")); origin != "" {
			c.Header(accessControlAllowOriginHeader, origin)
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
