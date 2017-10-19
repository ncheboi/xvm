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
	xvmName string
	xvmHome string
	osHome  string
	global  string
	local   string
	pwd     string
)

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

func argMin(min int) {
	if len(os.Args) < min {
		printGlobalFile("usage")
		os.Exit(1)
	}
}

func argMax(max int) {
	if len(os.Args) > max {
		printGlobalFile("usage")
		os.Exit(1)
	}
}

type Cmd func(args ...string)

func init() {
	if runtime.GOOS == "windows" {
		xvmName = "xvm"
		osHome = os.Getenv("USERPROFILE")
	} else {
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
		if x == xvmHome || isDir(xvmPath(x)) {
			local = xvmPath(x)
			break
		}

		x = filepath.Dir(x)
	}

	if local == "" {
		local = global
	}
}

func main() {
	argMin(2)

	cmd, ok := map[string]Cmd{
		"version": versionCmd,
		"usage":   usageCmd,
		"help":    helpCmd,
		"init":    initCmd,
		"which":   whichCmd,
		"status":  statusCmd,
		"remove":  removeCmd,
		"list":    listCmd,
		"pull":    pullCmd,
		"drop":    dropCmd,
		"set":     setCmd,
		"unset":   unsetCmd,
	}[os.Args[1]]

	if !ok {
		fail("Unknown command: %s", os.Args[1])
	}

	cmd()
}

func versionCmd(args ...string) {
	argMin(2)
	argMax(2)

	printGlobalFile("version")
}

func usageCmd(args ...string) {
	argMin(2)
	argMax(2)

	printGlobalFile("usage")
}

func helpCmd(args ...string) {
	argMin(2)
	argMax(2)

	printGlobalFile("README.md")
}

func initCmd(args ...string) {
	argMin(2)
	argMax(2)

	if local == xvmPath(pwd) {
		fail("Group already exists: %s", local)
	}

	mkdirAll(filepath.Join(xvmPath(pwd), "versions"))
}

func whichCmd(args ...string) {
	argMin(2)
	argMax(2)

	fmt.Println(local)
}

func statusCmd(args ...string) {
	argMin(2)
	argMax(3)

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
}

func removeCmd(args ...string) {
	argMin(2)
	argMax(2)

	removeAll(local)
}

func listCmd(args ...string) {
	argMin(3)
	argMax(4)

	p := pluginPath(os.Args[2])
	ensurePlugin(p)

	showAvailable := false

	if len(os.Args) > 3 {
		if os.Args[3] == "available" {
			showAvailable = true
		} else if os.Args[3] != "installed" {
			fail("Unknown argument: %s", os.Args[3])
		}
	}

	if showAvailable {
		printFile(filepath.Join(p, "available"))
		return
	}

	versions := dirNames(filepath.Join(p, "installed"))
	fmt.Println(strings.Join(versions, "\n"))
}

func pullCmd(args ...string) {
	argMin(3)
	argMax(4)

	p := pluginPath(os.Args[2])
	ensurePlugin(p)

	bin := filepath.Join(p, "bin", "pull")

	cmd := exec.Command(bin, os.Args[3])
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fail("%s", err)
	}
}

func dropCmd(args ...string) {
	argMin(4)
	argMax(4)

	p := pluginPath(os.Args[2])
	ensurePlugin(p)

	i := installedPath(p, os.Args[3])
	ensureInstalled(i)

	removeAll(i)
}

func setCmd(args ...string) {
	argMin(4)
	argMax(5)

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
}

func unsetCmd(args ...string) {
	argMin(3)
	argMax(4)

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
}
