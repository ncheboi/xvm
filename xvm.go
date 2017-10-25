package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/skotchpine/xvm/util"
)

// group specification options
const (
	Version = "0.0.2"

	Usage = `
xvm version
xvm usage
xvm help

xvm init
xvm which  [<pack>] [local|global]
xvm status [<pack>] [local|global]
xvm remove

xvm installed <pack>
xvm available <pack>
xvm stable    <pack>
xvm latest    <pack>

xvm set   <pack> <version> [local|global]
xvm unset <pack>           [local|global]

xvm pull <pack> <version>
xvm push <pack> <version>
xvm drop <pack> <version>

xvm config <pack> <version>

xvm alias   <pack> <version> <name>
xvm unalias <pack> <name>`

	StrGlobal = "global"
	StrLocal  = "local"

	StrPack      = "pack"
	StrPacks     = "packs"
	StrVersions  = "versions"
	StrInstalled = "installed"
	StrAvailable = "available"
	StrAliases   = "aliases"
	StrBin       = "bin"
	StrSplat     = "*"
)

var (
	globalDirPath, globalGroupPath string
	localDirPath, localGroupPath   string
	pwd                            string

	installedMap map[string][]string
	availableMap map[string][]string
	binMap       map[string]string
	aliasesMap   map[string]map[string]string

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
func FindGlobalGroup() (group, dir string) {
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
func FindLocalGroup() (group, dir string) {
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
func MapGroup(self, shared map[string]string, groupPath string) {
	versions, err := util.ReadMap(filepath.Join(groupPath, StrVersions))
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
func MapGroups(done chan bool) {
	MapGroup(localMap, currentMap, localGroupPath)
	MapGroup(globalMap, currentMap, globalGroupPath)
	done <- true
}

// Map the installed versions of all packages.
func MapInstalled(done chan bool) {
	defer func() { done <- true }()

	glob := filepath.Join(globalGroupPath, StrPacks, StrSplat, StrInstalled, StrSplat)
	list, err := filepath.Glob(glob)
	if err != nil {
		warn("Can not find installed versions")
		return
	}

	for _, path := range list {
		version := filepath.Base(path)
		pack := filepath.Base(filepath.Dir(filepath.Dir(path)))
		installedMap[pack] = append(installedMap[pack], version)
	}
}

// Map the executables for installed versions of all packages.
func MapBin(done chan bool) {
	defer func() { done <- true }()

	glob := filepath.Join(globalGroupPath, StrPacks, StrSplat, StrInstalled, StrSplat, StrBin, StrSplat)
	list, err := filepath.Glob(glob)
	if err != nil {
		warn("Can not find executable versions")
		return
	}

	for _, path := range list {
		dir := filepath.Dir
		pack := filepath.Base(dir(dir(dir(dir(path)))))
		binMap[pack] = path
	}
}

// Map the available versions of all packages.
func MapAvailable(done chan bool) {
	defer func() { done <- true }()

	glob := filepath.Join(globalGroupPath, StrPacks, StrSplat, StrAvailable)
	list, err := filepath.Glob(glob)
	if err != nil {
		warn("Can not find available versions")
		return
	}

	for _, path := range list {
		pack := filepath.Base(filepath.Dir(path))

		versions, err := util.ReadMap(path)
		if err != nil {
			warn("Can not find available versions for %s", pack)
			break
		}

		for version := range versions {
			availableMap[pack] = append(availableMap[pack], version)
		}
	}
}

// Map the aliases for all packages.
func MapAliases(done chan bool) {
	defer func() { done <- true }()

	glob := filepath.Join(globalGroupPath, StrPacks, StrSplat, StrAliases)
	list, err := filepath.Glob(glob)
	if err != nil {
		warn("Can not find aliases")
		return
	}

	for _, path := range list {
		pack := filepath.Base(filepath.Dir(path))

		aliases, err := util.ReadMap(path)
		if err != nil {
			warn("Can not find aliases for %s", pack)
			break
		}

		aliasesMap[pack] = aliases
	}
}

func ResolveAlias(pack, alias string) (concrete string) {
	if aliases, ok := aliasesMap[pack]; ok {
		if concrete, ok := aliases[alias]; ok {
			return concrete
		}
	}
	return alias
}

func WrapBin(bin string) {
	var pack, version string
	var ok bool

	if pack, ok = binMap[bin]; !ok {
		fail("Failed to find binary %s", bin)
	}
	if version, ok = currentMap[pack]; !ok {
		fail("No version set for package %s", pack)
	}

	path := filepath.Join(globalGroupPath, StrPacks, pack, StrInstalled, version, StrBin, bin)
	if util.NotExist(path) {
		fail("No executable %s for version %s of %s", bin, version, pack)
	}
	if util.Cmd(path) != nil {
		fail("")
	}
}

func init() {
	globalGroupPath, globalDirPath = FindGlobalGroup()
	localGroupPath, localDirPath = FindLocalGroup()

	done := make(chan bool)

	localMap = make(map[string]string)
	globalMap = make(map[string]string)
	currentMap = make(map[string]string)
	go MapGroups(done)

	availableMap = make(map[string][]string)
	go MapAvailable(done)

	installedMap = make(map[string][]string)
	go MapInstalled(done)

	binMap = make(map[string]string)
	go MapBin(done)

	aliasesMap = make(map[string]map[string]string)
	go MapAliases(done)

	for i := 0; i < 4; i++ {
		<-done
	}
}

func main() {
	// If the name of this file isn't xvm, but go, java, etc.,
	// find a relevant binary and execute it
	name := filepath.Base(os.Args[0])
	if name != "xvm"+OSExt {
		WrapBin(name)
	}

	if len(os.Args) < 2 {
		os.Args = append(os.Args, "usage")
	}

	switch os.Args[1] {
	case "version":
		fmt.Println(Version)
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
	default:
		fmt.Println(Usage)
	}
}

func argWrap(min, max int, fn func()) {
	n := len(os.Args)
	if (min > 0 && n < min) || (max > 0 && n > max) {
		fmt.Println(Usage)
	} else {
		fn()
	}
}

func initCmd() {
	if localGroupPath == pwd {
		fail("Group already exists")
	}
	path := filepath.Join(pwd, OSDir, StrVersions)
	if err := os.MkdirAll(path, util.PermPublic); err != nil {
		fail("")
	}
}

func whichCmd() {
	var group, pack string
	var ok bool

	for i := 2; i < len(os.Args); i++ {
		switch os.Args[i] {
		case StrGlobal, StrLocal:
			group = os.Args[i]
		default:
			pack = os.Args[i]
		}
	}

	if pack == "" {
		if group == StrGlobal {
			fmt.Println(globalDirPath)
		} else {
			fmt.Println(localDirPath)
		}
		return
	}

	if group != StrGlobal {
		if _, ok = localMap[pack]; ok {
			fmt.Println(localGroupPath)
			return
		}
	}

	if group != StrLocal {
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
		case StrGlobal, StrLocal:
			group = os.Args[i]
		default:
			pack = os.Args[i]
		}
	}

	var versions map[string]string
	switch group {
	case StrGlobal:
		versions = globalMap
	case StrLocal:
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
	if err := os.RemoveAll(filepath.Join(pwd, OSDir)); err != nil {
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
	if aliases, ok := aliasesMap[os.Args[2]]; ok {
		for alias := range aliases {
			fmt.Println(alias)
		}
	}
}

func stableCmd() {
	fmt.Println(ResolveAlias(os.Args[2], "stable"))
}

func latestCmd() {
	fmt.Println(ResolveAlias(os.Args[2], "latest"))
}

func setCmd() {
	pack := os.Args[2]
	version := ResolveAlias(pack, os.Args[3])

	base := localGroupPath
	if len(os.Args) == 5 {
		if os.Args[4] == StrGlobal {
			base = globalGroupPath
		} else if os.Args[4] != StrLocal {
			fmt.Println(Usage)
			os.Exit(1)
		}
	}

	if _, ok := installedMap[pack]; !ok {
		fail("Version %s of %s is not installed")
	}
	currentMap[pack] = version

	path := filepath.Join(base, StrVersions)
	if err := util.WriteMap(path, currentMap); err != nil {
		fail("Failed to save version")
	}
}

func unsetCmd() {
	pack := os.Args[2]

	base := localGroupPath
	if len(os.Args) == 4 {
		if os.Args[3] == StrGlobal {
			base = globalGroupPath
		} else if os.Args[3] != StrLocal {
			fmt.Println(Usage)
			os.Exit(1)
		}
	}

	if err := os.RemoveAll(filepath.Join(base, "versions", pack)); err != nil {
		fail(err.Error())
	}
}

func pullCmd() {
	pack := os.Args[2]
	version := ResolveAlias(pack, os.Args[3])

	var bin string
	if pack == StrPack {
		bin = filepath.Join(globalGroupPath, StrBin, "pull")
	} else {
		bin = filepath.Join(globalGroupPath, StrPacks, pack, StrInstalled, version, StrBin, "pull")
	}

	if err := util.Cmd(bin); err != nil {
		fail(err.Error())
	}
}

func dropCmd() {
	pack := os.Args[2]
	version := ResolveAlias(pack, os.Args[3])

	var path string
	if pack == StrPack {
		path = filepath.Join(globalGroupPath, StrPacks, version)
	} else {
		path = filepath.Join(globalGroupPath, StrPacks, pack, StrInstalled, version)
	}

	if err := os.RemoveAll(path); err != nil {
		fail(err.Error())
	}
}

func editCmd() {
	pack := os.Args[2]
	version := ResolveAlias(pack, os.Args[3])

	var path string
	if pack == StrPack {
		path = filepath.Join(globalGroupPath, StrPacks, version)
	} else {
		path = filepath.Join(globalGroupPath, StrPacks, pack, StrInstalled, version)
	}

	edit, ok := os.LookupEnv("EDITOR")
	if !ok || edit == "" {
		fail("Set EDITOR to edit config")
	}

	if err := util.Cmd(edit, path); err != nil {
		fail(err.Error())
	}
}

func authCmd() {
	fmt.Println("auth")
}

func pushCmd() {
	pack := os.Args[2]
	version := ResolveAlias(pack, os.Args[3])

	var bin string
	if pack == StrPack {
		bin = filepath.Join(globalGroupPath, StrBin, "pull")
	} else {
		bin = filepath.Join(globalGroupPath, StrPacks, pack, StrInstalled, version, StrBin, "pull")
	}

	if err := util.Cmd(bin); err != nil {
		fail(err.Error())
	}
}
