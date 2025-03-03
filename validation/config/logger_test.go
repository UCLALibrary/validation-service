//go:build unit

// Package config provides resources useful in the configuration of the validation service.
//
// This file contains a test to confirm the Zap middleware is working with the Echo server.
package config

import (
	"flag"
	"fmt"
	"github.com/UCLALibrary/validation-service/testflags"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// TestConfig loads the flags for the tests in the 'config' package.
func TestConfig(t *testing.T) {
	flag.Parse()
	fmt.Printf("%s's log level: %s\n", t.Name(), *testflags.LogLevel)
}

// TestZapLoggerMiddleware ensures the Zap middleware logs requests correctly
func TestZapLoggerMiddleware(t *testing.T) {
	logger := newTestLogger(t)
	server := echo.New()
	server.Use(ZapLoggerMiddleware(logger))

	// Define a test route
	server.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "test response")
	})

	// Create a test request
	request := httptest.NewRequest(http.MethodGet, "/test", nil)
	request.Header.Set("User-Agent", "TestAgent")
	recorder := httptest.NewRecorder()

	// Execute middleware and handler
	server.ServeHTTP(recorder, request)

	// Test the logger via the recorder
	assert.Equal(t, http.StatusOK, recorder.Code)
}

// newTestLogger gets a new logger to use in package's tests.
func newTestLogger(t *testing.T) *zap.Logger {
	config := zap.NewDevelopmentConfig()
	config.Level = zap.NewAtomicLevelAt(testflags.GetLogLevel())
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}

	logger, err := config.Build()
	if err != nil {
		t.Fatalf("Failed to build test logger: %v", err)
	}

	return logger
}
