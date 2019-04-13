package serverless

import (
	"fmt"
	"io"
	"os"
	"sync"
)

// Logger contains configuration for logging
type Logger struct {
	Target io.Writer
}

var (
	loggerInstance   *Logger
	loggerSetupNonce sync.Once
)

// GetLogger returns the shared logger instance.
func GetLogger() *Logger {
	loggerSetupNonce.Do(func() {
		loggerInstance = &Logger{Target: os.Stderr}
	})
	return loggerInstance
}

// Log writes the interfaces to the shared logger
func Log(a ...interface{}) {
	msg := fmt.Sprint(a...)
	LogMessage(msg)
}

// LogMessage writes the string to the shared logger
func LogMessage(msg string) {
	fmt.Fprintf(GetLogger().Target, "[SER] %s\n", msg)
}
