package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"./xvm"
	"./xvm/group"
)

var x = xvm.StartXvm()

func exitOn(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n",  err)
		os.Exit(1)
	}
}

func printDoc(name string) {
	path := filepath.Join(x.Path, name)

	contents, err := ioutil.ReadFile(path)
	exitOn(err)

	fmt.Print(string(contents))

	os.Exit(0)
}

func main() {
	if len(os.Args) < 2 {
		printDoc("usage")
	}

	switch os.Args[1] {
	case "version", "usage":
		printDoc(os.Args[1])
	case "help":
		printDoc("README.md")
	case "init", "which", "status", "remove", "set", "unset":
		groupCmd()
	case "list", "install", "uninstall", "update":
		pluginCmd()
	}
}

func groupCmd() {
	pwd, err := os.Getwd()
	exitOn(err)

	g := x.GetNearestGroup(pwd)

	switch os.Args[1] {
	case "init":
		if g.Path == filepath.Join(pwd, ".xvm") {
			exitOn(fmt.Errorf("A group already exists here"))
		}
		versions := filepath.Join(pwd, ".xvm", "versions")
		exitOn(os.MkdirAll(versions, 0755))

	case "which":
		fmt.Println(g.Path)

	case "status":
		if len(os.Args) > 2 {
			switch os.Args[2] {
			case "local":
			case "global":
				g = &group.Group{Path: x.Path}
			default:
				printDoc("usage")
			}
		}

		versions, err := g.GetVersions()
		exitOn(err)

		for plugin, version := range versions {
			fmt.Printf("%s %s\n", plugin, version)
		}
	
	case "remove":
		g.Remove()
	
	case "set":
		if len(os.Args) < 3 {
			printDoc("usage")
		}

		var plugin, version string

		if os.Args[2] == "local" {
			if len(os.Args) < 4 {
				printDoc("usage")
			}
			plugin = os.Args[3]
			version = os.Args[4]

		} else if os.Args[2] == "global" {
			if len(os.Args) < 4 {
				printDoc("usage")
			}
			g = x.GetNearestGroup(x.Path)
			plugin = os.Args[3]
			version = os.Args[4]

		} else {
			plugin = os.Args[2]
			version = os.Args[3]
		}

		exitOn(g.SetVersion(plugin, version))

	case "unset":
		if len(os.Args) < 3 {
			printDoc("usage")
		}

		var plugin string

		if os.Args[2] == "local" {
			if len(os.Args) < 3 {
				printDoc("usage")
			}
			plugin = os.Args[3]

		} else if os.Args[2] == "global" {
			if len(os.Args) < 3 {
				printDoc("usage")
			}
			g = x.GetNearestGroup(x.Path)
			plugin = os.Args[3]

		} else {
			plugin = os.Args[2]
		}

		exitOn(g.UnsetVersion(plugin))
	}
}

func pluginCmd() {
	if len(os.Args) < 3 {
		printDoc("usage")
	}

	if os.Args[2] == "plugin" {
		switch os.Args[1] {
		case "list":
			installed := true

			if len(os.Args) > 3 {
				switch os.Args[3] {
				case "installed":
				case "available":
					installed = false
				default:
					printDoc("usage")
				}
			}

			if installed {
				plugins, err := x.GetPluginsInstalled()
				exitOn(err)

				for _, plugin := range plugins {
					fmt.Println(plugin)
				}

			} else {
				plugins, err := x.GetPluginsAvailable()
				exitOn(err)

				for plugin := range plugins {
					fmt.Println(plugin)
				}
			}

		case "install":
			exitOn(fmt.Errorf("Plugin already installed"))
		case "uninstall":
			exitOn(fmt.Errorf("Plugin can't be uninstalled"))
		case "update":
			exitOn(fmt.Errorf("No updates for plugin"))
		}

		return
	}

	plugin, err := x.GetPlugin(os.Args[2])
	exitOn(err)

	switch os.Args[1] {
	case "list":
		installed := true

		if len(os.Args) > 3 {
			switch os.Args[3] {
			case "installed":
			case "available":
				installed = false
			default:
				printDoc("usage")
			}
		}

		var versions []string
		if installed {
			versions, err = plugin.GetInstalledVersions()
		} else {
			versions, err = plugin.GetAvailableVersions()
		}
		exitOn(err)

		for _, version := range versions {
			fmt.Println(version)
		}

	case "install":
		if len(os.Args) < 4 {
			printDoc("usage")
		}
		exitOn(plugin.Install(os.Args[1:]))

	case "uninstall":
		if len(os.Args) < 4 {
			printDoc("usage")
		}
		exitOn(plugin.Uninstall(os.Args[3]))

	case "update":
		exitOn(plugin.Update())
	}
}
