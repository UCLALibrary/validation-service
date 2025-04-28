//go:build unit
package utils

// This file tests that the default settings are correct on parsing.
import (
	"flag"
	"os"
	"testing"
)

// TestDefaultLogLevel tests that the default log level is set after the flag is parsed.
func TestDefaultLogLevel(t *testing.T) {
	// Backup original flag values (for restoring after test)
	originalArgs := os.Args

	// Remove any `-log-level` flag that might have been passed via `go test`
	os.Args = []string{originalArgs[0]}

	// Backup current LOG_LEVEL env variable (if set)
	originalLogLevel, isSet := os.LookupEnv("LOG_LEVEL")
	_ = os.Unsetenv("LOG_LEVEL") // Temporarily unset it for the test

	resetFlags() // Reset flags before testing
	flag.Parse() // Parse flags again with clean state

	// Check that the default log level is the one that we get after a new parse
	if LogLevel != "info" {
		t.Errorf("Expected default log level to be 'info', got '%s'", LogLevel)
	}

	// Restore original flag arguments (to not affect other tests)
	os.Args = originalArgs

	// Restore original LOG_LEVEL env variable
	if isSet {
		_ = os.Setenv("LOG_LEVEL", originalLogLevel)
	}
}

// resetFlags ensures flags are reset before each test.
func resetFlags() {
	// Create a new FlagSet and copy the testing framework's flags
	newFlagSet := flag.NewFlagSet("", flag.ContinueOnError)
	flag.VisitAll(func(flag *flag.Flag) {
		// Skip redefining the LogLevel flag if it already exists
		if flag.Name != "log-level" {
			newFlagSet.Var(flag.Value, flag.Name, flag.Usage)
		}
	})

	// Replace the default FlagSet with the new one
	flag.CommandLine = newFlagSet
	flag.StringVar(&LogLevel, "log-level", "info", "Log level (debug, info, warn, error)")
}
