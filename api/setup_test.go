//go:build unit
package api

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/UCLALibrary/validation-service/pkg/utils"
)

// TestMain loads the flags for the tests in the package.
func TestMain(main *testing.M) {
	flag.Parse()
	fmt.Printf("*** Package %s's log level: %s ***\n", utils.GetPackageName(), utils.LogLevel)
	os.Exit(main.Run())
}
