//go:build integration
package integration

// This file contains utilities for working with test loggers.
import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/UCLALibrary/validation-service/pkg/utils"
	"github.com/testcontainers/testcontainers-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// A logger to log our Docker container's logs
var logger *zap.Logger

// DockerLogConsumer implements the LogConsumer interface to capture container logs using the *zap.Logger.
type DockerLogConsumer struct {
	mutex sync.Mutex
}

// Accept relays the log entries from our Docker container to our test logger.
func (consumer *DockerLogConsumer) Accept(log testcontainers.Log) {
	consumer.mutex.Lock()
	defer consumer.mutex.Unlock()

	logger.Sugar().Debug(strings.TrimSpace(string(log.Content)))
}

// TcLogger wraps a Zap Logger for handling TestContainer messages.
type TcLogger struct {
	sugared *zap.SugaredLogger
}

// NewTcLogger creates a new zap.Logger wrapper that correctly reports line number.
func NewTcLogger(base *zap.Logger) *TcLogger {
	return &TcLogger{
		// Set caller skip to get the accurate line number from the unwrapped logger
		sugared: base.WithOptions(zap.AddCallerSkip(1)).Sugar(),
	}
}

// Printf prints a formatted message from the wrapped *zap.Logger
func (logger *TcLogger) Printf(format string, v ...interface{}) {
	message := strings.TrimSpace(fmt.Sprintf(format, v...)) // Removes extra LF at end from Docker output

	if shouldFilterLog(message) {
		return
	}

	logger.sugared.Info(message)
}

// getLogger builds a logger for use in relaying Docker container logs.
func getLogger(logLevel string) (*zap.Logger, error) {
	config := zap.NewDevelopmentConfig()

	// We don't need datestamps for running the tests, locally
	config.EncoderConfig.TimeKey = ""

	// Sets the current log level as the minimum logging level
	config.Level = zap.NewAtomicLevelAt(utils.GetLogLevel())

	// Configure the display of the log level
	config.EncoderConfig.LevelKey = "level"
	config.EncoderConfig.EncodeLevel = func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + level.CapitalString() + "]")
	}
	config.EncoderConfig.CallerKey = "caller"
	config.EncoderConfig.EncodeCaller = func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		_, file := filepath.Split(caller.File) // Extract just the file name from the full path
		enc.AppendString(fmt.Sprintf("%s:%d", file, caller.Line))
	}
	config.EncoderConfig.ConsoleSeparator = " "

	// Ensure logs are printed as plain text
	config.Encoding = "console"
	config.OutputPaths = []string{"stdout"}

	// Build and return the logger
	return config.Build()
}

// shouldFilterLog determines if the log message should be filtered or not.
func shouldFilterLog(message string) bool {
	return strings.Contains(message, "⏳") || strings.Contains(message, "✅")
}
