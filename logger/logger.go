package logger

import (
	"fmt"
	"log"
	"os"
)

// Loggable is the interface a logger must conform to
type Loggable interface {
	Output(int, string) error
}

// Logger is the default logger used by Serverless
var Logger Loggable

func init() {
	Logger = log.New(os.Stderr, "[SER] ", log.LstdFlags|log.Lmicroseconds|log.LUTC)
}

// Print writes the interfaces to the shared logger
func Print(a ...interface{}) {
	msg := fmt.Sprint(a...)
	Logger.Output(2, msg)
}

// Printf writes the formated log message to the shared logger.
func Printf(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	Logger.Output(2, msg)
}
