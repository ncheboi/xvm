package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// group specification options
const (
	OptGlobal = "global"
	OptLocal  = "local"

	PackName     = "pack"
	PacksName    = "packs"
	VersionsName = "versions"
)

var (
	globalDirPath, globalGroupPath string
	localDirPath, localGroupPath   string
	pwd                            string

	installedMap map[string][]string
	availableMap map[string][]string
	binMap       map[string]string

	localMap   map[string]string
	globalMap  map[string]string
	currentMap map[string]string
)

func warn(msg string, etc ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", etc...)
}

func fail(msg string, etc ...interface{}) {
	warn(msg, etc...)
	os.Exit(1)
}

// Use XVMPATH as the global group. If XVMPATH is not set, the global
// group resolve by appending the default directory name to the user's home.
func findGlobalGroup() (group, dir string) {
	var ok bool
	group, ok = os.LookupEnv("XVMPATH")
	if !ok || group == "" {
		home, ok := os.LookupEnv(OSHome)
		if !ok || home == "" {
			fail("Could not resolve XVMPATH. Either set XVMPATH or %s", OSHome)
		}
		group = filepath.Join(home, OSDir)
	}
	return group, filepath.Dir(group)
}

// Set the nearest group. If none exist between the working directory
// and the root, use the global group.
func findLocalGroup() (group, dir string) {
	group = globalGroupPath

	var err error
	pwd, err = os.Getwd()
	if err != nil {
		warn("Failed to get working directory")
		goto returnLocalGroup
	}

	// Move from the current directory to the root; stop before crossing the global path.
	for x := pwd; x != "/" && x != globalDirPath; x = filepath.Dir(x) {
		// If a group is found (xvm directory exists), it is the local group.
		info, err := os.Stat(filepath.Join(x, OSDir))
		if err == nil && info.IsDir() {
			group = x
			break
		}
	}

returnLocalGroup:
	return group, filepath.Dir(group)
}

// Add a group's versions to self and shared version maps.
// Do not overwrite existing entries in the shared map.
func mapGroup(self, shared map[string]string, path string, done chan bool) {
	defer func() { done <- true }()

	versions, err := ReadConfig(path)
	if err != nil {
		warn(err.Error())
		return
	}

	for pack, version := range versions {
		// Write each version to this group's map.
		self[pack] = version

		// Write each version to the shared map only if no entry exists.
		if _, ok := shared[pack]; !ok {
			shared[pack] = version
		}
	}
}

// Map the versions of the local and global group. Give the local group precedence.
func mapGroups(done chan bool) {
	mapGroup(localMap, currentMap, localGroupPath)
	mapGroup(globalMap, currentMap, globalGroupPath)
	done <- true
}

// Map the installed versions of all packages.
//
// NOTE: this is nearly an O(n^3) implementation. There must be a less runtime-
//       intensive way to keep this available.
func mapInstalled(done chan bool) {
	defer func() { done <- true }()

	i := join(globalGroupPath, "installed")

	packs, err := dirnames(i)
	if err != nil {
		warn("Failed to search installed packages")
		a <- true
		return
	}

	for _, pack := range packs {
		j := join(i, pack, "installed")

		versions, err := dirnames(j)
		if err != nil {
			if !os.IsNotExist(err) {
				warn("Failed to list versions for %s", pack)
			}
			continue
		}
		installedMap[pack] = versions

		for _, version := range versions {
			k := join(j, version, "bin")

			bins, err := dirnames(k)
			if err != nil {
				if !os.IsNotExist(err) {
					warn("Failed to list binaires for version %s of %s", version, pack)
				}
				continue
			}

			for _, bin := range bins {
				if _, ok := binMap[bin]; !ok {
					binMap[bin] = pack
				}
			}
		}
	}
}

// Map the available versions of all packages
func mapAvailable(done chan bool) {
	defer func() { done <- true }()

	i := join(globalGroupPath, "available")

	packs, err := dirnames(i)
	if err == nil {
		availableMap[PackName] = packs
	} else {
		warn("Failed to list available packages")
	}

	i = join(globalGroupPath, "installed")

	packs, err = dirnames(i)
	if err != nil {
		warn("Failed to search installed packages")
		b <- true
	}

	for _, pack := range packs {
		j := join(i, pack, "available")

		versions, err := dirnames(j)
		if err != nil {
			if !os.IsNotExist(err) {
				warn("Failed to list available versions for %s", pack)
			}
			continue
		}
		availableMap[pack] = versions
	}
}

func wrapBinary(bin string) {
	var pack, version string
	var ok bool

	if pack, ok = binMap[bin]; !ok {
		fail("Failed to find binary %s", bin)
	}
	if version, ok = currentMap[pack]; !ok {
		fail("No version set for package %s", pack)
	}

	path := join(globalGroupPath, "installed", pack, "installed", version, "bin", bin)
	if notexist(path) {
		fail("No executable %s for version %s of %s", bin, version, pack)
	}
	if cmd(path) != nil {
		fail("")
	}
}

func init() {
	globalGroupPath, globalDirPath = findGlobalGroup()
	localGroupPath, localDirPath = findLocalGroup()

	done := make(chan bool)

	localMap = make(map[string]string)
	globalMap = make(map[string]string)
	currentMap = make(map[string]string)
	go mapGroups(done)

	availableMap = make(map[string][]string)
	go mapAvailable(done)

	installedMap = make(map[string][]string)
	binMap = make(map[string]string)
	go mapInstalled(done)

	for i = 0; i < 3; i++ {
		<-done
	}
}

func main() {
	// If the name of this file isn't xvm, but go, java, etc.,
	// find a relevant binary and execute it
	name := filepath.Base(os.Args[0])
	if name != "xvm"+OSExt {
		wrapBinary(name)
	}

	if len(os.Args) < 2 {
		printfile(globalGroupPath, "usage")
	}

	switch os.Args[1] {
	case "init":
		argWrap(2, 2, initCmd)
	case "which":
		argWrap(2, 4, whichCmd)
	case "current":
		argWrap(2, 4, currentCmd)
	case "remove":
		argWrap(2, 2, removeCmd)
	case "installed":
		argWrap(3, 3, installedCmd)
	case "available":
		argWrap(3, 3, availableCmd)
	case "stable":
		argWrap(3, 3, stableCmd)
	case "latest":
		argWrap(3, 3, latestCmd)
	case "set":
		argWrap(4, 5, setCmd)
	case "unset":
		argWrap(3, 4, unsetCmd)
	case "pull":
		argWrap(4, 4, pullCmd)
	case "drop":
		argWrap(4, 4, dropCmd)
	case "edit":
		argWrap(3, 3, editCmd)
	case "auth":
		argWrap(3, 3, authCmd)
	case "push":
		argWrap(3, 3, pushCmd)

	case "version":
		printfile(globalGroupPath, "version")
	case "help":
		printfile(globalGroupPath, "readme")
	default:
		printfile(globalGroupPath, "usage")
	}
}

func argWrap(min, max int, fn func()) {
	n := len(os.Args)
	if (min > 0 && n < min) || (max > 0 && n > max) {
		printfile(globalGroupPath, "usage")
	}
	fn()
}

func initCmd() {
	if localGroupPath == pwd {
		fail("Group already exists")
	}
	if mkdir(pwd, OSDir, "versions") != nil {
		fail("")
	}
}

func whichCmd() {
	var group, pack string
	var ok bool

	for i := 2; i < len(os.Args); i++ {
		switch os.Args[i] {
		case OptGlobal, OptLocal:
			group = os.Args[i]
		default:
			pack = os.Args[i]
		}
	}

	if pack == "" {
		if group == OptGlobal {
			fmt.Println(globalDirPath)
		} else {
			fmt.Println(localDirPath)
		}
		return
	}

	if group != OptGlobal {
		if _, ok = localMap[pack]; ok {
			fmt.Println(localGroupPath)
			return
		}
	}

	if group != OptLocal {
		if _, ok = globalMap[pack]; ok {
			fmt.Println(globalGroupPath)
			return
		}
	}
}

func currentCmd() {
	var group, pack string
	for i := 2; i < len(os.Args); i++ {
		switch os.Args[i] {
		case OptGlobal, OptLocal:
			group = os.Args[i]
		default:
			pack = os.Args[i]
		}
	}

	var versions map[string]string
	switch group {
	case OptGlobal:
		versions = globalMap
	case OptLocal:
		versions = localMap
	default:
		versions = currentMap
	}

	if pack != "" {
		if version, ok := versions[pack]; ok {
			fmt.Println(version)
		}
		return
	}

	for pack, version := range versions {
		fmt.Printf("%s %s\n", pack, version)
	}
}

func removeCmd() {
	if localGroupPath == globalGroupPath {
		fail("Cannot remove global group")
	}
	if localGroupPath != pwd {
		fail("Group does not exist")
	}
	if err := rm(pwd, OSDir); err != nil {
		fail(err.Error())
	}
}

func installedCmd() {
	if versions, ok := installedMap[os.Args[2]]; ok {
		for _, version := range versions {
			fmt.Println(version)
		}
	}
}

func availableCmd() {
	if versions, ok := availableMap[os.Args[2]]; ok {
		for _, version := range versions {
			fmt.Println(version)
		}
	}
}

func stableCmd() {
	if version, err := alias(os.Args[2], "stable"); err == nil {
		fmt.Println(version)
	}
}

func latestCmd() {
	if version, err := alias(os.Args[2], "latest"); err == nil {
		fmt.Println(version)
	}
}

func setCmd() {
	pack, version := os.Args[2], os.Args[3]
	if v, err := alias(pack, version); err == nil {
		version = v
	}

	base := localGroupPath
	if len(os.Args) == 5 {
		if os.Args[4] == OptGlobal {
			base = globalGroupPath
		} else if os.Args[4] != OptLocal {
			printfile(globalGroupPath, "usage")
		}
	}

	if _, ok := installedMap[pack]; !ok {
		fail("Version %s of %s is not installed")
	}

	if writeline(version, base, "versions", pack) != nil {
		fail("Failed to save version")
	}
}

func unsetCmd() {
	pack := os.Args[2]

	base := localGroupPath
	if len(os.Args) == 4 {
		if os.Args[3] == OptGlobal {
			base = globalGroupPath
		} else if os.Args[3] != OptLocal {
			printfile(globalGroupPath, "usage")
		}
	}

	if err := rm(base, "versions", pack); err != nil {
		fail(err.Error())
	}
}

func pullCmd() {
	pack, version := os.Args[2], os.Args[3]
	if v, err := alias(pack, version); err == nil {
		version = v
	}

	var bin string
	if pack == PackName {
		bin = join(globalGroupPath, "bin", "pull")
	} else {
		bin = join(globalGroupPath, "installed", pack, "installed", version, "bin", "pull")
	}

	if cmd(bin) != nil {
		fail("")
	}
}

func dropCmd() {
	pack, version := os.Args[2], os.Args[3]
	if v, err := alias(pack, version); err == nil {
		version = v
	}

	var path string
	if pack == PackName {
		path = join(globalGroupPath, "installed", version)
	} else {
		path = join(globalGroupPath, "installed", pack, "installed", version)
	}

	if rm(path) != nil {
		fail("")
	}
}

func editCmd() {
	pack, version := os.Args[2], os.Args[3]
	if v, err := alias(pack, version); err == nil {
		version = v
	}

	var path string
	if pack == PackName {
		path = join(globalGroupPath, "installed", version)
	} else {
		path = join(globalGroupPath, "installed", pack, "installed", version)
	}

	edit, ok := os.LookupEnv("EDITOR")
	if !ok || edit == "" {
		fail("Set EDITOR to edit config")
	}

	if cmd(edit, path) != nil {
		fail("")
	}
}

func authCmd() {
	fmt.Println("auth")
}

func pushCmd() {
	pack, version := os.Args[2], os.Args[3]
	if v, err := alias(pack, version); err == nil {
		version = v
	}

	var bin string
	if pack == PackName {
		bin = join(globalGroupPath, "bin", "pull")
	} else {
		bin = join(globalGroupPath, "installed", pack, "installed", version, "bin", "pull")
	}

	if cmd(bin) != nil {
		fail("")
	}
}
