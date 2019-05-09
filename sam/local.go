package sam

import "os"

// IsLocal returns a boolean true if the application is running in a SAM local environment
func IsLocal() bool {
	return os.Getenv("AWS_SAM_LOCAL") == "true"
}
