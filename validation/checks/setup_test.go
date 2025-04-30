//go:build unit

package checks

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
	fmt.Printf("*** Package %s's HOST_DIR: %s ***\n", utils.GetPackageName(), utils.HostDir)

	// We handle the error of it being missing in the check that uses it
	_ = os.Setenv("HOST_DIR", utils.HostDir)
	os.Exit(main.Run())
}
