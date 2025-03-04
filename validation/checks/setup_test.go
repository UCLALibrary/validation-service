//go:build unit

// Package checks consists of individual validation checks.
//
// This file sets up the testing environment for the tests.
package checks

import (
	"flag"
	"fmt"
	"github.com/UCLALibrary/validation-service/testflags"
	"testing"
)

// TestCheck loads the flags for the tests in the 'checks' package.
func TestCheck(t *testing.T) {
	flag.Parse()
	fmt.Printf("%s's log level: %s\n", t.Name(), *testflags.LogLevel)
}
