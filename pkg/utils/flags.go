//go:build unit || integration

// Package utils provides tools to help with testing.
package utils

// This file defines LogLevel and a has a convenience function for constructing levels.
import (
	"flag"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var ServiceName string   // The name of the Web service we're building, running, and testing
var LogLevel string      // The desired log level at which the Web service should be run
var HostDir string       // The location of the mounted host directory (contains resource files)
var KakaduVersion string // The version of Kakadu that's being used (this may also be empty if Kakadu isn't used)
var BuildArch string     // The system architecture that's going to run the Kakadu build (default is usually fine)

// init initializes the flags that have been defined
func init() {
	flag.StringVar(&ServiceName, "service-name", "service", "Name of service being tested")
	flag.StringVar(&LogLevel, "log-level", "info", "Log level (debug, info, warn, error)")
	flag.StringVar(&HostDir, "host-dir", "", "HOST_DIR env variable that is copied into test-container")
	flag.StringVar(&KakaduVersion, "kakadu-version", "", "The version of Kakadu being used")
	flag.StringVar(&BuildArch, "arch", "x86-64", "The system architecture of the Kakadu build")
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
