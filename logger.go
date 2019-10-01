package serverless

import (
	"fmt"
	"log"
	"os"
)

// Logger is the default logger used by Serverless
var Logger *log.Logger

func init() {
	Logger = log.New(os.Stderr, "[SER] ", log.LstdFlags|log.Lmicroseconds|log.LUTC)
}

// GetLogger returns the shared logger instance.
func GetLogger() *log.Logger {
	return Logger
}

// Log writes the interfaces to the shared logger
func Log(a ...interface{}) {
	msg := fmt.Sprint(a...)
	Logger.Print(msg)
}

// Logf writes the formated log message to the shared logger.
func Logf(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	Logger.Print(msg)
}
