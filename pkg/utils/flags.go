//go:build unit || integration

// Package utils provides tools to help with testing.
//
// This file defines LogLevel and a has a convenience function for constructing levels.
package utils

import (
	"flag"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var ServiceName string // The name of the Web service we're building, running, and testing
var LogLevel string    // The desired log level at which the Web service should be run
var HostDir string     // The location of the mounted host directory (contains resource files)

// init initializes the flags that have been defined
func init() {
	flag.StringVar(&ServiceName, "service-name", "service", "Name of service being tested")
	flag.StringVar(&LogLevel, "log-level", "info", "Log level (debug, info, warn, error)")
	flag.StringVar(&HostDir, "host-dir", "", "HOST_DIR env variable that is copied into test-container")
}

// GetLogLevel gets the current log level as a zapcore.Level.
func GetLogLevel() zapcore.Level {
	switch LogLevel {
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
