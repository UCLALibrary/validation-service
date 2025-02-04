// Package testflags provides flags that can be used in testing.
//
// This file defines LogLevel and a has a convenience function for constructing levels.
package testflags

import (
	"flag"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LogLevel is defined as a global test flag
var LogLevel = flag.String("log-level", "info", "Log level (debug, info, warn, error)")

// GetLogLevel converts a string version of a Zap log level to zapcore.Level.
func GetLogLevel(level *string) zapcore.Level {
	// Use the default level if a nil is passed in
	if level == nil {
		return zap.InfoLevel
	}

	// Otherwise, dereference pointer and return corresponding Zap level
	switch *level {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	default:
		return zap.InfoLevel
	}
}
