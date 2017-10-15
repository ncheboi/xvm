package main

import (
	"os"
	"io/ioutil"
	"fmt"
	"log"
	"path/filepath"

	"./utils"
	"./group"
	"./plugin"
	"./delegate"
)

// Print a file in $XVMPATH/doc to Stdout.
func printDoc(name string) {
	fullpath := filepath.Join(utils.XvmPath(), "doc", name)
	buf, err := ioutil.ReadFile(fullpath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(buf))
}

// Print documentation from $XVMPATH/doc,
// delegate other arguments to subpackages.
func main() {
	args := os.Args
	if (len(args) < 2) {
		printDoc("usage")
		return
	}

	switch args[1] {
	case "version", "usage", "help":
		printDoc(args[1])
	case "group", "set", "unset":
		group.Main(args[2:])
	case "plugin", "uninstall":
		plugin.Main(args[2:])
	case "list", "install":
		delegate.Main(args[2:])
	default:
		printDoc("usage")
	}
}
