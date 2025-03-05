//go:build unit

// Package config provides resources useful in the configuration of the validation service.
//
// This file sets up the testing environment for the package.
package config

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
	fmt.Printf("*** Package %s's log level: %s ***\n", utils.GetPackageName(), *utils.LogLevel)
	os.Exit(main.Run())
}
