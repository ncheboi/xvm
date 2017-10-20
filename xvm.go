package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

var (
	osBinName string
	platform string
	xvmName  string
	xvmHome  string
	osHome   string
	global   string
	local    string
	pwd      string
)

type Cmd struct {
	min, max int
	fn       func()
}

func (cmd *Cmd) run() {
	if cmd.min > 0 && len(os.Args) < cmd.min {
		printGlobalFile("usage")
		os.Exit(1)
	}

	if cmd.min > 0 && len(os.Args) > cmd.max {
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

func interpretVersion(p, version string) version string {
	if version != "stable" && version != "latest" {
		if isFile(filepath.Join(p, version)) {
			version = readLineOnce(path)
		}
	}
}

func readLineOnce(path string) string {
	file, err := os.Open(path)
	failIf(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	return scanner.Text()
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
		osBinName = "xvm.exe"
		platform = "windows"
		xvmName = "xvm"
		osHome = os.Getenv("USERPROFILE")
	} else {
		osBinName = "xvm"
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
	"version": {2, 2, // xvm version
		func() {
			printGlobalFile("version")
		},
	},

	"usage": {2, 2, // xvm usage
		func() {
			printGlobalFile("usage")
		},
	},

	"help": {2, 2, // xvm help
		func() {
			printGlobalFile("readme")
		},
	},

	"init": {2, 2, // xvm init
		func() {
			if local == xvmPath(pwd) {
				fail("Group already exists: %s", local)
			}
			mkdirAll(filepath.Join(xvmPath(pwd), "versions"))
		},
	},

	"which": {2, 3, // xvm which [<plugin>]
		func() {
			if len(os.Args) == 3 {
				path := filepath.Join(local, "versions", os.Args[2])
				if !isFile(path) {
					fmt.Println(global)
					return
				}
			}

			fmt.Println(local)
		},
	},

	"status": {2, 4, // xvm status [local|global] [<plugin>]
		func() {
			which := local

			path := filepath.Join(which, "versions")
			if !isDir(path) {
				mkdirAll(path)
			}

			for i := 2; i < len(os.Args); i++ {
				if os.Args[i] == "global" {
					which = global
				} else if os.Args[i] != "local" {
					whence := filepath.Join(path, os.Args[i])
					if isFile(whence) {
						printFile(whence)
					}
					return
				}
			}

			for _, plugin := range dirNames(path) {
				contents, err := ioutil.ReadFile(path)
				failIf(err)

				fmt.Printf("%s %s", plugin, string(contents))
			}
		},
	},

	"remove": {2, 2, // xvm remove
		func() {
			removeAll(local)
		},
	},

	"show": {3, 4, // xvm show <plugin> [installed|available|stable|latest]
		func() {
			plugin := os.args[2]
			p := pluginPath(plugin)
			ensurePlugin(p)

			var kind string

			if len(os.Args) > 3 {
				switch os.Args[3] {
				case "stable", "latest":
					path := filepath.Join(p, os.Args[3])
					if isFile(path) {
						printFile(path)
					}
					return

				case "installed":
					kind = "installed"
				case "available":
					kind = "available"

				default:
					fail("Unknown argument: %s", os.Args[3])
				}
			}

			path := filepath.Join(p, kind)
			if isDir(path) {
				for _, whence := range dirNames(path) {
					fmt.Println(whence)
				}
			}
		},
	},

	"pull": {3, 4, // xvm pull <plugin> <version>|stable|latest
		func() {
			plugin := os.Args[2]
			p := pluginPath(plugin)
			ensurePlugin(p)

			version := interpretVersion(os.Args[3])
			destDir := installedPath(p, version)
			if !isDir(destDir) {
				mkdirAll(destDir)
			}

			path := filepath.Join(p, "available", version)
			content, err := ioutil.ReadFile(path)
			failIf(err)

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

	"drop": {4, 4, // xvm drop <plugin> <version>|stable|latest
		func() {
			plugin := os.Args[2]
			p := pluginPath(plugin)
			ensurePlugin(p)

			version := interpretVersion(os.Args[3])
			i := installedPath(p, version)
			ensureInstalled(i)

			removeAll(i)
		},
	},

	"set": {4, 5, // xvm set [local|global] <plugin> <version>|stable|latest
		func() {
			which := local
			plugin := os.Args[2]
			version := os.Args[3]

			if len(os.Args) == 5 {
				if os.Args[2] == "global" {
					which = global
				} else if os.Args[2] != "local" {
					fail("Unknown argument: %s", os.Args[2])
				}

				plugin = os.Args[3]
				version = os.Args[4]
			}

			if plugin == "plugin" {
				fmt.Println("Plugin cannot be set")
			}

			p := pluginPath(plugin)
			ensurePlugin(p)

			version = interpretVersion(version)
			i := installedPath(p, version)
			ensureInstalled(i)

			path := filepath.Join(which, "versions", plugin)
			writeFile(path, filepath.Base(i)+"\n")
		},
	},

	"unset": {3, 4, // xvm unset [local|global] <plugin> <version>
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

	"shim": {2, 2, // xvm shim
		func() {
		}
	},
}

func main() {
	if len(os.Args) < 2 {
		printGlobalFile("usage")
		os.Exit(1)
	}

	if os.Args[0] != osBinName {
		exec()
		return
	}

	cmd, ok := cmds[os.Args[1]]
	if !ok {
		fail("Unknown command: %s", os.Args[1])
	}
	cmd.run()
}

func shim() {
}

func exec() {
	if len(os.Args) < 2 {
		fail("")
	}

	args := os.Args[2:]
	cmd := exec.Command(bin, [][]byte(args))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fail("%s", err)
	}
}
