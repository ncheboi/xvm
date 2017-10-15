package main

import (
	"os"
	"io/ioutil"
	"fmt"
	"log"
	"path"
	"path/filepath"
)

func dumpFile(name string) {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	fullpath := filepath.Join(path.Dir(pwd), "doc", name)
	buf, err := ioutil.ReadFile(fullpath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(buf))
}

func group(args []string) {}
func plugin(args []string) {}
func delegate(args []string) {}

func main() {
	args := os.Args
	if (len(args) < 2) {
		return;
	}

	switch args[1] {
	case "version":
		dumpFile("version")
	case "usage":
		dumpFile("usage")
	case "help":
		dumpFile("help")
	case "group", "set", "unset":
		group(args)
	case "plugin", "uninstall":
		plugin(args)
	case "list", "install":
		delegate(args)
	}
}
