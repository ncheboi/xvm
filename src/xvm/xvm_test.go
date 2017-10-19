package xvm_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"."
)

func TestXvmPath(t *testing.T) {
	os.Setenv("XVMPATH", "/xvmpath")
	if xvm := xvm.StartXvm(); xvm.Path != "/xvmpath" {
		t.Fail()
	}
}

func TestGetNearestGroup(t *testing.T) {
	root := filepath.Join(os.TempDir(), "xvmtest")
	os.RemoveAll(root); defer os.RemoveAll(root)

	testPath := func(names... string) string {
		path := append([]string{root}, names...)
		return filepath.Join(path...)
	}

	paths := []string{
		testPath("home", ".xvm"),
		testPath("proj", "x", ".xvm"),
		testPath("proj", "x", "lib", "y", ".xvm"),
		testPath("etc", "etc"),
	}
	for _, path := range paths {
		if err := os.MkdirAll(path, 0755); err != nil {
			t.Errorf("Failed to make xvm dir: %s", path)
		}
	}

	os.Setenv("XVMPATH", filepath.Join(root, "home", ".xvm"))
	xvm := xvm.StartXvm()

	tests := map[string]string {
		testPath("home"):                  testPath("home", ".xvm"),
		testPath("proj"):                  testPath("home", ".xvm"),
		testPath("etc", "etc"):            testPath("home", ".xvm"),
		testPath("proj", "x"):             testPath("proj", "x", ".xvm"),
		testPath("proj", "x", "lib"):      testPath("proj", "x", ".xvm"),
		testPath("proj", "x", "lib", "y"): testPath("proj", "x", "lib", "y", ".xvm"),
	}

	for given, expected := range tests {
		group := xvm.GetNearestGroup(given)
		if group.Path != expected {
			t.Errorf("Expected nearest group of %s to be %s, but found %s", given, expected, group.Path)
		}
	}
}

func TestGetPluginsAvailable(t *testing.T) {
	root := filepath.Join(os.TempDir(), "xvmtest")
	os.RemoveAll(root); os.MkdirAll(root, 0755)
	defer os.RemoveAll(root)

	path := filepath.Join(root, "available")
	file, err := os.Create(path)
	if err != nil {
		t.Errorf("Failed to create %s: %s", path, err)
	}
	defer file.Close()

	expected := map[string]string{
		"a": "1",
		"b": "2",
	}
	for plugin, version := range expected {
		file.WriteString(plugin + " " + version + "\n")
	}

	os.Setenv("XVMPATH", root)
	xvm := xvm.StartXvm()

	available, err := xvm.GetPluginsAvailable()
	if err != nil {
		t.Error(err)
	}
	for expectedPlugin, expectedVersion := range expected {
		found := false

		for actualPlugin, actualVersion := range available {
			if actualPlugin == expectedPlugin {
				found = true
				break

				if actualVersion != expectedVersion {
					t.Error("Failed to report correct URL")
				}
			}
		}

		if !found {
			t.Errorf("Failed to report available plugin: %s", expectedPlugin)
		}
	}
}

func TestGetPluginsInstalled(t *testing.T) {
	root := filepath.Join(os.TempDir(), "xvmtest")
	os.RemoveAll(root); defer os.RemoveAll(root)

	testPath := func(names... string) string {
		path := append([]string{root}, names...)
		return filepath.Join(path...)
	}

	expected := []string{
		testPath("installed", "v1"),
		testPath("installed", "v2"),
	}
	for _, path := range expected {
		if err := os.MkdirAll(path, 0755); err != nil {
			t.Errorf("Failed to make xvm dir: %s", path)
		}
	}

	os.Setenv("XVMPATH", testPath())
	xvm := xvm.StartXvm()

	installed, err := xvm.GetPluginsInstalled()
	if err != nil {
		t.Error(err)
	}
	for _, path := range expected {
		expected := filepath.Base(path)
		found := false

		for _, actual := range installed {
			if actual == expected {
				found = true
				break
			}
		}

		if !found {
			t.Error(fmt.Errorf("Failed to report installed plugin: %s", expected))
		}
	}
}

func TestGetPlugin(t *testing.T) {
	root := filepath.Join(os.TempDir(), "xvmtest")
	os.RemoveAll(root); os.MkdirAll(root, 0755)
	defer os.RemoveAll(root)

	path := filepath.Join(root, "available")
	file, err := os.Create(path)
	if err != nil {
		t.Errorf("Failed to create %s: %s", path, err)
	}
	defer file.Close()

	expected := map[string]string{
		"a": "1",
		"b": "2",
	}
	for plugin, version := range expected {
		file.WriteString(plugin + " " + version + "\n")
	}

	os.Setenv("XVMPATH", root)
	xvm := xvm.StartXvm()

	for name, url := range expected {
		plugin, err := xvm.GetPlugin(name)
		if err != nil {
			t.Error(err)
		}
		if plugin.Url != url {
			t.Errorf("Failed to assign plugin URL")
		}
	}
}
