// +build windows

package main

// Platform-specific filesystem defaults.
const (
	OSExt  = ".exe"        // windows binaries need an extension; go compiles to *.exe
	OSDir  = "xvm"         // name of directory for local groups
	OSHome = "USERPROFILE" // path of default global group
)
