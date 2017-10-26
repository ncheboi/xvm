package main_test

import (
	"os"
	"path/filepath"
	"testing"

	xvm "github.com/skotchpine/xvm"
)

var (
	root = filepath.Join(os.TempDir(), "xvm-main-test")
)

func TestGlobalGroup(t *testing.T) {
	expectedDir := filepath.Join(root, "HOME")
	expectedGroup := filepath.Join(expectedDir, "XVM")
	os.Setenv("XVMPATH", expectedGroup)

	actualGroup, actualDir := xvm.FindGlobalGroup()
	if actualGroup != expectedGroup {
		t.Errorf("Expected %s, got %s", expectedGroup, actualGroup)
	}
	if actualDir != expectedDir {
		t.Errorf("Expected %s, got %s", expectedDir, actualDir)
	}

	expectedDir = filepath.Join(root, xvm.OSHome, xvm.OSHome)
	expectedGroup = filepath.Join(expectedDir, xvm.OSDir)
	os.Unsetenv("XVMPATH")
	os.Setenv(xvm.OSHome, expectedDir)

	actualGroup, actualDir = xvm.FindGlobalGroup()
	if actualGroup != expectedGroup {
		t.Errorf("Expected %s, got %s", expectedGroup, actualGroup)
	}
	if actualDir != expectedDir {
		t.Errorf("Expected %s, got %s", expectedDir, actualDir)
	}

	os.Unsetenv("XVMPATH")
	os.Unsetenv(xvm.OSHome)
}

func TestLocalGroup(t *testing.T) {
	mkdir := func(elem ...string) string {
		path := filepath.Join(elem...)
		if err := os.MkdirAll(path, 0777); err != nil {
			t.Error(err)
		}
		path, _ = filepath.EvalSymlinks(path)
		return path
	}

	near := mkdir(root, "users")
	group := mkdir(root, "users", "project")
	groupxvm := mkdir(root, "users", "project", xvm.OSDir)
	subdir := mkdir(root, "users", "project", "subproject")
	xvmpath := mkdir(root, "users", "home", "xvm")

	defer func() {
		if err := os.RemoveAll(filepath.Join(root, "users")); err != nil {
			t.Error(err)
		}
	}()

	os.Setenv("XVMPATH", xvmpath)
	xvm.GlobalGroupPath, xvm.GlobalDirPath = xvm.FindGlobalGroup()

	tests := []struct{ pwd, group string }{
		{xvmpath, xvmpath}, // XVMPATH, itself
		{group, groupxvm},  // a group, itself
		{subdir, groupxvm}, // a subdir, its group
		{near, xvmpath},    // near root, XVMPATH
	}

	for _, expected := range tests {
		if err := os.Chdir(expected.pwd); err != nil {
			t.Error(err)
		}

		group, _ := xvm.FindLocalGroup()
		if group != expected.group {
			t.Errorf("Expected %s, got %s", expected.group, group)
		}
	}
}

func TestMapGroup(t *testing.T) {
	t.Skip()
}

func TestMapInstalled(t *testing.T) {
	t.Skip()
}

func TestMapBin(t *testing.T) {
	t.Skip()
}

func TestMapAvailable(t *testing.T) {
	t.Skip()
}

func TestMapAliases(t *testing.T) {
	t.Skip()
}

func TestResolveAlias(t *testing.T) {
	t.Skip()
}

func TestWrapBin(t *testing.T) {
	t.Skip()
}
