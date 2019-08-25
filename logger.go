package serverless

import (
	"fmt"
	"log"
	"os"
	"sync"
)

var (
	loggerInstance   *log.Logger
	loggerSetupNonce sync.Once
)

// GetLogger returns the shared logger instance.
func GetLogger() *log.Logger {
	loggerSetupNonce.Do(func() {
		loggerInstance = log.New(os.Stderr, "[SER] ", log.LstdFlags|log.Lmicroseconds|log.LUTC)
	})
	return loggerInstance
}

// Log writes the interfaces to the shared logger
func Log(a ...interface{}) {
	msg := fmt.Sprint(a...)
	GetLogger().Print(msg)
}

// Logf writes the formated log message to the shared logger.
func Logf(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	GetLogger().Print(msg)
}
