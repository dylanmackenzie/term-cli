package logger

import (
	"log"

	"github.com/wsxiaoys/terminal/color"
)

// Wrapper for log.Printf
func Log(format string, args ...interface{}) {
	c := color.Colorize("y!")
	log.Printf(c+format+color.ResetCode, args...)
}

// Wrapper for log.Fatalf
func Fatal(format string, args ...interface{}) {
	c := color.Colorize("r")
	log.Fatalf(c+format+color.ResetCode, args...)
}
