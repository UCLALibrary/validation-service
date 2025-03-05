//go:build unit

// Package utils provides utilities that can be used in testing.
//
// This file sets up the package's testing environment.
package utils

import (
	"flag"
	"fmt"
	"os"
	"testing"
)

// TestMain loads the flags for the tests in the package.
func TestMain(main *testing.M) {
	flag.Parse()
	fmt.Printf("*** Package %s's log level: %s ***\n", GetPackageName(), *LogLevel)
	os.Exit(main.Run())
}
