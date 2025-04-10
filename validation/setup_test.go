//go:build unit

// Package validation provides tools to validate CSV data.
//
// This file sets up the package's testing environment.
package validation

import (
	"flag"
	"fmt"
	"github.com/UCLALibrary/validation-service/pkg/utils"
	"os"
	"testing"
)

// TestMain loads the flags for the tests in the package.
func TestMain(main *testing.M) {
	flag.Parse()
	fmt.Printf("*** Package %s's log level: %s ***\n", utils.GetPackageName(), utils.LogLevel)
	os.Exit(main.Run())
}
