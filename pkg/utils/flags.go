//go:build unit

// Package utils provides tools to help with testing.
//
// This file defines LogLevel and a has a convenience function for constructing levels.
package utils

import (
	"flag"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LogLevel is defined as a global test flag
var LogLevel = flag.String("log-level", "info", "Log level (debug, info, warn, error)")

// GetLogLevel gets the current log level as a zapcore.Level.
func GetLogLevel() zapcore.Level {
	switch *LogLevel {
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
