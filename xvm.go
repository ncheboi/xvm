package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	platform string
	xvmName string
	xvmHome string
	osHome  string
	global  string
	local   string
	pwd     string
)

type Cmd struct {
	min, max int
	fn       func()
}

func (cmd *Cmd) run() {
	if len(os.Args) < cmd.min {
		printGlobalFile("usage")
		os.Exit(1)
	}

	if len(os.Args) > cmd.max {
		printGlobalFile("usage")
		os.Exit(1)
	}

	cmd.fn()
}

func fail(msg string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", a...)
	os.Exit(1)
}

func failIf(err error) {
	if err != nil {
		fail(err.Error())
	}
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	failIf(err)
	return info.IsDir()
}

func isFile(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	failIf(err)
	return !info.IsDir() && info.Mode().IsRegular()
}

func removeAll(path string) {
	failIf(os.RemoveAll(path))
}

func mkdirAll(path string) {
	failIf(os.MkdirAll(path, 0755))
}

func writeFile(path string, content string) {
	failIf(ioutil.WriteFile(path, []byte(content), 0755))
}

func dirNames(path string) []string {
	dir, err := os.Open(path)
	failIf(err)

	names, err := dir.Readdirnames(0)
	failIf(err)

	return names
}

func xvmPath(path string) string {
	return filepath.Join(path, xvmName)
}

func pluginPath(name string) string {
	if name == "plugin" {
		return global
	}
	return filepath.Join(global, "installed", name)
}

func ensurePlugin(p string) {
	if !isDir(p) {
		fail("Not installed: %s", filepath.Base(p))
	}
}

func installedPath(p, name string) string {
	return filepath.Join(p, "installed", name)
}

func ensureInstalled(i string) {
	if !isDir(i) {
		fail("Not installed: %s", filepath.Base(i))
	}
}

func printGlobalFile(name string) {
	printFile(filepath.Join(global, name))
}

func printLocalFile(name string) {
	printFile(filepath.Join(local, name))
}

func printFile(path string) {
	contents, err := ioutil.ReadFile(path)
	failIf(err)
	fmt.Print(string(contents))
}

func init() {
	if runtime.GOOS == "windows" {
		platform = "windows"
		xvmName = "xvm"
		osHome = os.Getenv("USERPROFILE")
	} else {
		platform = "unix"
		xvmName = ".xvm"
		osHome = os.Getenv("HOME")
	}
	xvmHome = osHome

	var isSet bool
	global, isSet = os.LookupEnv("XVMPATH")
	if isSet {
		xvmHome = filepath.Dir(global)
	} else {
		global = xvmPath(osHome)
	}

	pwd, err := os.Getwd()
	if err != nil {
		fail("Failed to get working directory")
	}

	x := pwd
	for x != "/" {
		if x == xvmHome {
			local = global
			break
		}

		if isDir(xvmPath(x)) {
			local = xvmPath(x)
			break
		}

		x = filepath.Dir(x)
	}

	if local == "" {
		local = global
	}
}

var cmds = map[string]Cmd{
	"version": {2, 2,
		func() {
			printGlobalFile("version")
		},
	},

	"usage": {2, 2,
		func() {
			printGlobalFile("README.md")
		},
	},

	"help": {2, 2,
		func() {
			printGlobalFile("help")
		},
	},

	"init": {2, 2,
		func() {
			if local == xvmPath(pwd) {
				fail("Group already exists: %s", local)
			}
			mkdirAll(filepath.Join(xvmPath(pwd), "versions"))
		},
	},

	"which": {2, 2,
		func() {
			fmt.Println(local)
		},
	},

	"status": {2, 3,
		func() {
			which := local

			if len(os.Args) > 2 {
				if os.Args[2] == "global" {
					which = global
				} else if os.Args[2] != "local" {
					fail("Unknown argument: %s", os.Args[2])
				}
			}

			versionsPath := filepath.Join(which, "versions")
			if !isDir(versionsPath) {
				mkdirAll(versionsPath)
			}
			plugins := dirNames(versionsPath)

			for _, plugin := range plugins {
				printFile(filepath.Join(versionsPath, plugin))
			}
		},
	},

	"remove": {2, 2,
		func() {
			removeAll(local)
		},
	},

	"list": {3, 4,
		func() {
			p := pluginPath(os.Args[2])
			ensurePlugin(p)

			which := "installed"

			if len(os.Args) > 3 {
				if os.Args[3] == "available" {
					which = "available"
				} else if os.Args[3] != "installed" {
					fail("Unknown argument: %s", os.Args[3])
				}
			}

			path := filepath.Join(p, which)
			if isDir(path) {
				versions := dirNames(path)
				fmt.Println(strings.Join(versions, "\n"))
			}
		},
	},

	"pull": {3, 4,
		func() {
			plugin := os.Args[2]
			version := os.Args[3]

			p := pluginPath(plugin)
			ensurePlugin(p)

			destDir := installedPath(p, version)
			if !isDir(destDir) {
				mkdirAll(destDir)
			}

			contentPath := filepath.Join(p, "available", version)
			rawContent, err := ioutil.ReadFile(contentPath)
			failIf(err)

			content := strings.TrimRight(string(rawContent), "\n")

			bin := filepath.Join(p, platform, "pull")

			os.Setenv("XVM_PULL_DESTDIR", destDir)
			os.Setenv("XVM_PULL_VERSION", version)
			os.Setenv("XVM_PULL_CONTENT", content)

			cmd := exec.Command(bin)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {
				fail("%s", err)
			}
		},
	},

	"drop": {4, 4,
		func() {
			p := pluginPath(os.Args[2])
			ensurePlugin(p)

			i := installedPath(p, os.Args[3])
			ensureInstalled(i)

			removeAll(i)
		},
	},

	"set": {4, 5,
		func() {
			which := local
			plugin := os.Args[2]
			version := os.Args[3]

			if len(os.Args) > 4 {
				if os.Args[2] == "global" {
					which = global
				} else if os.Args[2] != "local" {
					fail("Unknown argument: %s", os.Args[2])
				}

				plugin = os.Args[3]
				version = os.Args[4]
			}

			p := pluginPath(plugin)
			ensurePlugin(p)

			i := installedPath(p, version)
			ensureInstalled(i)

			path := filepath.Join(which, "versions", plugin)
			writeFile(path, filepath.Base(i)+"\n")
		},
	},

	"unset": {3, 4,
		func() {
			which := local
			plugin := os.Args[3]

			if len(os.Args) > 3 {
				if os.Args[2] == "global" {
					which = global
				} else if os.Args[2] != "local" {
					fail("Unknown argument: %s", os.Args[2])
				}
			} else {
				plugin = os.Args[2]
			}

			ensurePlugin(pluginPath(plugin))
			removeAll(filepath.Join(which, "versions", plugin))
		},
	},
}

func main() {
	if len(os.Args) < 2 {
		printGlobalFile("usage")
		os.Exit(1)
	}

	cmd, ok := cmds[os.Args[1]]
	if !ok {
		fail("Unknown command: %s", os.Args[1])
	}
	cmd.run()
}
