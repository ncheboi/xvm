package xvm

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"os"
	"runtime"
	"strings"

	"./plugin"
	"./group"
)

type Xvm struct {
	Path     string // $XVMPATH, $USERPROFILE/xvm or $HOME/.xvm
	Platform string // windows or unix
}

// Creates a new Xvm instance.
func StartXvm() *Xvm {
	xvm := new(Xvm)

	if runtime.GOOS == "windows" {
		xvm.Platform = "windows"
	} else {
		xvm.Platform = "unix"
	}

	env, isSet := os.LookupEnv("XVMPATH")
	if isSet {
		xvm.Path = env
	} else {
		switch xvm.Platform {
		case "windows":
			xvm.Path = filepath.Join(os.Getenv("USERPROFILE"), "xvm")
		case "unix":
			xvm.Path = filepath.Join(os.Getenv("HOME"), ".xvm")
		}
	}

	return xvm
}

// Returns a group for the first directory containing a .xvm directory between
// the path and the root or a group for the $XVMPATH.
func (x *Xvm) GetNearestGroup(path string) *group.Group {
	for path != "/" && path != x.Path {
		groupPath := filepath.Join(path, ".xvm")

		info, err := os.Stat(groupPath)
		if err == nil && info.IsDir() {
			return &group.Group{Path: groupPath}
		}

		path = filepath.Dir(path)
	}

	return &group.Group{Path: x.Path}
}

// Get a map of plugins and their update urls from the $XVMPATH/available.
func (x *Xvm) GetPluginsAvailable() (map[string]string, error) {
	path := filepath.Join(x.Path, "available")

	plugins := make(map[string]string)

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines[:len(lines) - 1] {
		tokens := strings.Split(line, " ")
		if len(tokens) != 2 {
			return nil, fmt.Errorf("%s is ill-formatted", path)
		}

		plugins[tokens[0]] = tokens[1]
	}

	return plugins, nil
}

// Get a list of all installed plugins from $XVMPATH/installed.
func (x *Xvm) GetPluginsInstalled() ([]string, error) {
	path := filepath.Join(x.Path, "installed")
	dir, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return dir.Readdirnames(0)
}

// Creates a Plugin for a given name.
func (x *Xvm) GetPlugin(name string) (*plugin.Plugin, error) {
	available, err := x.GetPluginsAvailable()
	if err != nil {
		return nil, err
	}

	url, ok := available[name]
	if !ok {
		return nil, fmt.Errorf("Plugin not found: %s", name)
	}

	p := &plugin.Plugin{
		Path: filepath.Join(x.Path, "installed", name),
		Bin:  filepath.Join(x.Path, "bin." + x.Platform, name),
		Url:  url,
	}

	return p, nil
}
