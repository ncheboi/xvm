package utils

import (
	"fmt"
	"os"
	"runtime"
	"path/filepath"
)

// Gets an absolute path to Xvm's root directory. Use XVMPATH if set.
// Returns an error if the file doesn't exist or isn't a directory.
// If the file could not be validated, returns the path and an error.
func GetXvmPath() (string, error) {
	path, isSet := os.LookupEnv("XVMPATH")
	if !isSet {
		var home string

		if runtime.GOOS == "windows" {
			home = os.Getenv("USERPROFILE")
		} else {
			home = os.Getenv("HOME")
		}

		path = filepath.Join(home, ".xvm")
	}

	info, err := os.Stat(path)
	if info == nil {
		return path, fmt.Errorf("XVMPATH was resolved to %s, but cannot be validated", path)
	}
	if os.IsNotExist(err) || !info.IsDir() {
		return "", fmt.Errorf("XVMPATH was resolved to %s, which is not a directory", path)
	}

	return path, nil
}

// Gets an absolute path for a subdirectory of Xvm's root, using XVMPATH if set.
// Returns an error if the file doesn't exist or isn't a directory.
// If the file could not be validated, returns the path and an error.
func GetXvmSubDir(path string) (string, error) {
	root, err := GetXvmPath()
	if err != nil {
		return "", err
	}
	absPath := filepath.Join(root, path)

	info, err := os.Stat(absPath)
	if info == nil {
		return path, fmt.Errorf("%s cannot be validated", absPath)
	}
	if os.IsNotExist(err) || !info.IsDir() {
		return "", fmt.Errorf("%s is not a directory", absPath)
	}

	return absPath, nil
}
