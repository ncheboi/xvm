// +build !windows

package main

// Platform-specific filesystem defaults.
const (
	OSExt  = ""     // unix binaries need no extensions
	OSDir  = ".xvm" // name of hidden directory for local groups
	OSHome = "HOME" // path of default global group
)
