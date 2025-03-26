//go:build integration

// Package integration holds the project's integration tests.
//
// This file contains utilities for working with test loggers.
package integration

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// A logger to log our Docker container's logs
var logger *zap.Logger

// DockerLogConsumer implements the LogConsumer interface to capture container logs.
type DockerLogConsumer struct {
	mutex sync.Mutex
}

// FilteredLogger filters specific TestContainer messages.
type FilteredLogger struct {
	original *log.Logger
}

// getLogger builds a logger for use in relaying Docker container logs.
func getLogger() (*zap.Logger, error) {
	config := zap.NewDevelopmentConfig()

	config.EncoderConfig.TimeKey = ""   // Remove Time Field
	config.EncoderConfig.LevelKey = ""  // Remove Level Field
	config.EncoderConfig.CallerKey = "" // Remove Caller Field
	config.EncoderConfig.MessageKey = "message"

	// Simplify the label of the logging output
	config.EncoderConfig.EncodeTime = func(time time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[DOCKER]")
	}

	// Get rid of line breaks in the output
	config.EncoderConfig.ConsoleSeparator = ""

	// Ensure logs are printed as plain text
	config.Encoding = "console"
	config.OutputPaths = []string{"stdout"}

	// Build and return the logger
	logger, _ := config.Build()
	//noinspection GoUnhandledErrorResult
	defer logger.Sync()

	return logger, nil
}

// shouldFilterLog determines if the log message should be filtered or not.
func shouldFilterLog(message string) bool {
	return strings.Contains(message, "⏳") || strings.Contains(message, "✅")
}

// Accept relays the log entries from our Docker container to our test logger.
func (consumer *DockerLogConsumer) Accept(log testcontainers.Log) {
	consumer.mutex.Lock()
	defer consumer.mutex.Unlock()

	// Output our container's logs to our test logger
	logger.Debug(string(log.Content))
}

// Print prints the log messages we're not filtering.
func (logger *FilteredLogger) Print(log ...interface{}) {
	message := fmt.Sprint(log...)

	// Suppress logs that contain don't provide useful information
	if shouldFilterLog(message) {
		return
	}

	// Allow other logs to pass through
	logger.original.Print(log...)
}

// Println prints the log messages we're not filtering.
func (logger *FilteredLogger) Println(log ...interface{}) {
	message := strings.TrimSpace(strings.Join(strings.Fields(strings.Trim(fmt.Sprint(log...), "[]")), " "))

	// Suppress logs that contain don't provide useful information
	if shouldFilterLog(message) {
		return
	}

	// Allow all other logs to pass through
	logger.original.Println(log...)
}

// Printf prints the log messages we're not filtering.
func (logger *FilteredLogger) Printf(format string, log ...interface{}) {
	message := fmt.Sprintf(format, log...)

	// Suppress logs that contain don't provide useful information
	if shouldFilterLog(message) {
		return
	}

	// Allow all other logs to pass through
	logger.original.Printf(format, log...)
}
