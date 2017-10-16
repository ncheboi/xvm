package main

import (
	"os"
	"fmt"

	_"./utils"
	"./docs"
	"./group"
	"./plugin"
	"./delegate"
)

// Print an error to Stderr and exit with a given code.
func exitWithError(err error, code int) {
	fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
	os.Exit(code)
}

// Controls program exit. Delegates to subcommands.
func main() {
	args := os.Args
	if (len(args) < 2) {
		if err := docs.Print("usage"); err != nil {
			exitWithError(err, 1)
		}
		return
	}

	switch args[1] {
	case "version", "usage", "help":
		if err := docs.Print(args[1]); err != nil {
			exitWithError(err, 1)
		}

	case "group", "set", "unset":
		if err := group.Main(args); err != nil {
			exitWithError(err, 1)
		}

	case "plugin", "uninstall":
		if err := plugin.Main(args); err != nil {
			exitWithError(err, 1)
		}

	case "list", "install":
		if err := delegate.Main(args); err != nil {
			exitWithError(err, 1)
		}

	default:
		if err := docs.Print("usage"); err != nil {
			exitWithError(err, 1)
		}
	}
}
