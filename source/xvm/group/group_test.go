package group_test

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"
	"."
)

func MakeVersionsDir(t *testing.T) string {
	tmp := os.TempDir()

	dir := filepath.Join(tmp, "group", "versions")
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		t.Errorf("Failed to make directories %s", dir)
	}

	return dir
}

func TestGetVersions(t *testing.T) {
	dir := MakeVersionsDir(t)
	g := &group.Group{path.Dir(dir)}

	expected := map[string]string{
		"a": "1",
		"b": "2",
	}

	for plugin, version := range expected {
		path := filepath.Join(dir, plugin)
		err := ioutil.WriteFile(path, []byte(version + "\n"), os.ModePerm)
		if err != nil {
			t.Errorf("Failed to write version file %s", path)
		}
	}

	actual, err := g.GetVersions()
	if err != nil {
		t.Error(err)
	}

	for plugin, version := range expected {
		if actual[plugin] != version {
			t.Fail()
		}
	}
}

func TestSetUnsetVersion(t *testing.T) {
	dir := MakeVersionsDir(t)
	g := &group.Group{path.Dir(dir)}

	err := g.SetVersion("a", "1")
	if err != nil {
		t.Error(err)
	}

	ver, err := g.GetVersion("a")
	if err != nil {
		t.Error(err)
	}

	if ver != "1" {
		t.Fail()
	}

	err = g.UnsetVersion("a")
	if err != nil {
		t.Error(err)
	}

	_, err = g.GetVersion("a")
	if err == nil {
		t.Fail()
	}
}

func TestRemove(t *testing.T) {
	dir := MakeVersionsDir(t)
	g := &group.Group{path.Dir(dir)}

	err := g.Remove()
	if err != nil {
		t.Error(err)
	}

	_, err = os.Stat(dir)
	if !os.IsNotExist(err) {
		t.Fail()
	}
}
