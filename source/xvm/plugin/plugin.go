package plugin

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Plugin struct {
	Path string // The root of the plugin
	Bin  string // Path of this platform's install execuatable
	Url  string // Points to a dist.zip file available over http
}

// Gets the versions published in plugin/available.
func (p *Plugin) GetAvailableVersions() ([]string, error) {
	path := filepath.Join(p.Path, "available")

	versions, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return strings.Split(string(versions), "\n"), nil
}

// Gets all versions installed to plugin/installed and not removed.
func (p *Plugin) GetInstalledVersions() ([]string, error) {
	path := filepath.Join(p.Path, "installed")

	dir, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return dir.Readdirnames(0)
}

// Delegates all arguments to this platform's install executable.
func (p *Plugin) Install(args []string) error {
	cmd := exec.Command(p.Bin, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Removes an installed version from plugin/installed.
func (p *Plugin) Uninstall(version string) error {
	path := filepath.Join(p.Path, "installed", version)
	return os.RemoveAll(path)
}

// Removes the entire plugin.
func (p *Plugin) Remove() error {
	return os.RemoveAll(p.Path)
}
