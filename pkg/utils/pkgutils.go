//go:build unit || integration

package utils

import (
	"path"
	"runtime"
)

// GetPackageName gets the current package name dynamically.
func GetPackageName() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Base(path.Dir(filename))
}
