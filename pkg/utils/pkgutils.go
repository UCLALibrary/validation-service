//go:build unit || integration

package utils

// This file contains testing utils related to packages and packaging.
import (
	"path"
	"runtime"
)

// GetPackageName gets the current package name dynamically.
func GetPackageName() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Base(path.Dir(filename))
}
