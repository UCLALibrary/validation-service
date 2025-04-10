//go:build unit || integration

// Package utils has structures and utilities useful for running tests.
//
// This file contains testing utils related to packages and packaging.
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
