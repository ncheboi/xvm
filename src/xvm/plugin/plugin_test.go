package plugin_test

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"testing"
	"strings"

	"."
)

func MakePlugin(t *testing.T) *plugin.Plugin {
	tmp := os.TempDir()
	subdir := filepath.Join(tmp, "plugin", "installed")
	if err := os.MkdirAll(subdir, os.ModePerm); err != nil {
		t.Errorf("Failed to make directories %s", subdir)
	}

	echo, err := exec.LookPath("echo")
	if err != nil {
		t.Error("Failed to find echo command")
	}

	return &plugin.Plugin{
		Path: path.Dir(subdir),
		Bin: echo,
		Url: "localhost:9876",
	}
}

func RemovePlugin(t *testing.T, p *plugin.Plugin) {
	if err := os.RemoveAll(p.Path); err != nil {
		t.Error(err)
	}
}

func TestGetAvailableVersions(t *testing.T) {
	p := MakePlugin(t)
	defer RemovePlugin(t, p)

	expected := []string{"10", "9", "8", "7"}

	path := filepath.Join(p.Path, "available")
	err := ioutil.WriteFile(path, []byte(strings.Join(expected, "\n")), os.ModePerm)
	if err != nil {
		t.Errorf("Failed to make available file %s", path)
	}

	actual, err := p.GetAvailableVersions()
	if err != nil {
		t.Error(err)
	}
	for i, a := range actual {
		if a != expected[i] {
			t.Fail()
		}
	}
}

func TestGetInstalledVersions(t *testing.T) {
	p := MakePlugin(t)
	defer RemovePlugin(t, p)

	expected := []string{"1", "2"}

	for _, e := range expected {
		path := filepath.Join(p.Path, "installed", e)
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			t.Errorf("Failed to make installed version %s", path)
		}
	}

	actual, err := p.GetInstalledVersions()
	if err != nil {
		t.Error(err)
	}
	for _, a := range actual {
		found := false

		for _, e := range expected {
			if a == e {
				found = true
				break
			}
		}

		if found == false {
			t.Fail()
		}
	}
}

func TestUninstall(t *testing.T) {
	p := MakePlugin(t)
	defer RemovePlugin(t, p)

	path := filepath.Join(p.Path, "installed", "1")
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		t.Errorf("Failed to make installed version %s", path)
	}

	if err := p.Uninstall("1"); err != nil {
		t.Error(err)
	}
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		t.Fail()
	}
}

func TestRemove(t *testing.T) {
	p := MakePlugin(t)
	defer RemovePlugin(t, p)

	if err := p.Remove(); err != nil {
		t.Error(err)
	}
	_, err := os.Stat(p.Path)
	if !os.IsNotExist(err) {
		t.Fail()
	}
}
