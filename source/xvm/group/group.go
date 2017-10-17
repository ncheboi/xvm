package group

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Group struct {
	Path string
}

// Gets a map of plugins and their versions from the versions subdirectory.
// Returns an error if a plugin's version file isn't as expected.
func (g *Group) GetVersions() (map[string]string, error) {
	path := filepath.Join(g.Path, "versions")

	dir, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	plugins, err := dir.Readdirnames(0)
	if err != nil {
		return nil, err
	}

	versions := make(map[string]string)

	for _, plugin := range plugins {
		version, err := g.GetVersion(plugin)
		if err != nil {
			return nil, err
		}

		versions[plugin] = version
	}

	return versions, nil
}

// Gets a version from a plugin's file in the versions subdirectory.
// Returns an error if the plugin's version file isn't as expected.
func (g *Group) GetVersion(plugin string) (string, error) {
	path := filepath.Join(g.Path, "versions", plugin)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	n := bytes.Count(content, []byte{'\n'})
	if n > 1 {
		return "", fmt.Errorf("%s is ill-formatted", path)
	}

	return strings.Trim(string(content), "\n"), nil
}

// Writes a plugin file with a version in the versions subdirectory.
// Truncates the plugin file if it exists.
func (g *Group) SetVersion(plugin string, version string) error {
	path := filepath.Join(g.Path, "versions", plugin)
	return ioutil.WriteFile(path, []byte(version + "\n"), os.ModePerm)
}

// Removes a plugin file form the versions subdirectory.
func (g *Group) UnsetVersion(plugin string) error {
	path := filepath.Join(g.Path, "versions", plugin)
	return os.RemoveAll(path)
}

// Removes the group's directory and everything in it.
func (g *Group) Remove() error {
	return os.RemoveAll(g.Path)
}
