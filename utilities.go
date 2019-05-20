package serverless

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/maddiesch/serverless/sam"
)

// IsDebug returns a boolean if the current gin mode is not release
func IsDebug() bool {
	return sam.IsLocal() || gin.Mode() == gin.ReleaseMode
}

// GetenvDefault fetches an env variable. If it's blank it will return the default.
func GetenvDefault(key string, value string) string {
	val := os.Getenv(key)
	if len(val) == 0 {
		return value
	}
	return val
}
