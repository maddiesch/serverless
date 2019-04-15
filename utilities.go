package serverless

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// IsLocal returns a boolean if the app is running in the SAM local env
func IsLocal() bool {
	return strings.Compare(os.Getenv("AWS_SAM_LOCAL"), "true") == 0
}

// IsDebug returns a boolean if the current gin mode is not release
func IsDebug() bool {
	return IsLocal() || strings.Compare(gin.Mode(), gin.ReleaseMode) != 0
}

// DefaultEnv fetches an env variable. If it's blank it will return the default.
func DefaultEnv(key string, value string) string {
	val := os.Getenv(key)
	if len(val) == 0 {
		return value
	}
	return val
}
