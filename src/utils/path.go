package utils

import (
	"os"
	"runtime"
	"path/filepath"
)

// Get an absolute path to Xvm's root directory. Use XVMPATH if set.
func XvmPath() string {
	path, isSet := os.LookupEnv("XVMPATH")
	if isSet {
		return path
	}

	if runtime.GOOS == "windows" {
		path = os.Getenv("USERPROFILE")
	} else {
		path = os.Getenv("HOME")
	}

	return filepath.Join(path, ".xvm")
}
